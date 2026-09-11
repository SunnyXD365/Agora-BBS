package workflow

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"agora-backend/internal/llm"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	temporalworkflow "go.temporal.io/sdk/workflow"
)

type CoolingInput struct {
	Kind           string
	ID             int64
	EndsAt         time.Time
	RequiresReview bool
}

type CoolingStarter interface {
	StartCooling(context.Context, CoolingInput) error
}

type TemporalStarter struct {
	client    client.Client
	taskQueue string
}

func NewTemporalStarter(temporalClient client.Client, taskQueue string) *TemporalStarter {
	return &TemporalStarter{client: temporalClient, taskQueue: taskQueue}
}

func (s *TemporalStarter) StartCooling(ctx context.Context, input CoolingInput) error {
	workflowID := fmt.Sprintf("cooling-%s-%d-%d", input.Kind, input.ID, input.EndsAt.UnixNano())
	_, err := s.client.ExecuteWorkflow(ctx, client.StartWorkflowOptions{ID: workflowID, TaskQueue: s.taskQueue}, ContentCoolingWorkflow, input)
	return err
}

func ContentCoolingWorkflow(ctx temporalworkflow.Context, input CoolingInput) (string, error) {
	if delay := input.EndsAt.Sub(temporalworkflow.Now(ctx)); delay > 0 {
		if err := temporalworkflow.Sleep(ctx, delay); err != nil {
			return "", err
		}
	}
	ctx = temporalworkflow.WithActivityOptions(ctx, temporalworkflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy:         &temporal.RetryPolicy{InitialInterval: time.Second, BackoffCoefficient: 2, MaximumInterval: 30 * time.Second, MaximumAttempts: 5},
	})
	var result PublishResult
	if err := temporalworkflow.ExecuteActivity(ctx, "PublishContent", input).Get(ctx, &result); err != nil {
		return "", err
	}
	if input.Kind == "post" && result.Status == "published" && result.PostCount >= 20 && (result.PostCount == 20 || result.PostCount%10 == 0) {
		_ = temporalworkflow.ExecuteActivity(ctx, "ClusterTopic", result.TopicID).Get(ctx, nil)
	}
	if input.Kind == "topic" && result.Status == "pending_review" {
		child := temporalworkflow.WithChildOptions(ctx, temporalworkflow.ChildWorkflowOptions{WorkflowID: fmt.Sprintf("blind-review-topic-%d", input.ID)})
		if err := temporalworkflow.ExecuteChildWorkflow(child, BlindReviewWorkflow, BlindReviewInput{SubjectType: "topic", SubjectID: input.ID}).Get(child, nil); err != nil {
			return result.Status, err
		}
	}
	return result.Status, nil
}

type PublishResult struct {
	Status    string
	TopicID   int64
	PostCount int
}

type Activities struct {
	DB  *sql.DB
	LLM llm.Client
}

func (a *Activities) PublishContent(ctx context.Context, input CoolingInput) (PublishResult, error) {
	activity.RecordHeartbeat(ctx, input.Kind, input.ID)
	tx, err := a.DB.BeginTx(ctx, nil)
	if err != nil {
		return PublishResult{}, err
	}
	defer tx.Rollback()
	status := "published"
	resultData := PublishResult{Status: status}
	if input.Kind == "topic" {
		if input.RequiresReview {
			status = "pending_review"
		}
		result, err := tx.ExecContext(ctx, `UPDATE topics SET status = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND status = 'cooling' AND cooling_ends_at <= CURRENT_TIMESTAMP`, input.ID, status)
		if err != nil {
			return PublishResult{}, err
		}
		rows, _ := result.RowsAffected()
		if rows == 0 {
			status = "superseded"
		}
	} else if input.Kind == "post" {
		var topicID int64
		err := tx.QueryRowContext(ctx, `UPDATE posts SET status = 'published', updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND status = 'cooling' AND cooling_ends_at <= CURRENT_TIMESTAMP RETURNING topic_id`, input.ID).Scan(&topicID)
		if errors.Is(err, sql.ErrNoRows) {
			status = "superseded"
		} else if err != nil {
			return PublishResult{}, err
		} else {
			if _, err := tx.ExecContext(ctx, `UPDATE topics SET post_count = post_count + 1 WHERE id = $1`, topicID); err != nil {
				return PublishResult{}, err
			}
			resultData.TopicID = topicID
			if err := tx.QueryRowContext(ctx, `SELECT post_count FROM topics WHERE id = $1`, topicID).Scan(&resultData.PostCount); err != nil {
				return PublishResult{}, err
			}
		}
	} else {
		return PublishResult{}, fmt.Errorf("unknown content kind %q", input.Kind)
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO governance_events(event_key, event_type, aggregate_type, aggregate_id, payload, status, processed_at)
		VALUES ($1, 'cooling_completed', $2, $3, jsonb_build_object('result', $4::text), 'completed', CURRENT_TIMESTAMP)
		ON CONFLICT(event_key) DO NOTHING`, fmt.Sprintf("cooling:%s:%d:%d", input.Kind, input.ID, input.EndsAt.UnixNano()), input.Kind, input.ID, status)
	if err != nil {
		return PublishResult{}, err
	}
	resultData.Status = status
	return resultData, tx.Commit()
}

func registerActivities(w interface {
	RegisterActivityWithOptions(interface{}, activity.RegisterOptions)
}, activities *Activities) {
	w.RegisterActivityWithOptions(activities.PublishContent, activity.RegisterOptions{Name: "PublishContent"})
}
