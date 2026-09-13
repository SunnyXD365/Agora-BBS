package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"agora-backend/internal/dao"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	temporalworkflow "go.temporal.io/sdk/workflow"
)

type BlindReviewInput struct {
	SubjectType string
	SubjectID   int64
}
type BlindReviewStarter interface {
	StartBlindReview(context.Context, string, int64) error
	StartReviewAudit(context.Context, int64) error
	StartFinalizeReview(context.Context, int64) error
}

func (s *TemporalStarter) StartBlindReview(ctx context.Context, subjectType string, subjectID int64) error {
	_, err := s.client.ExecuteWorkflow(ctx, client.StartWorkflowOptions{ID: fmt.Sprintf("blind-review-%s-%d-%d", subjectType, subjectID, time.Now().UnixNano()), TaskQueue: s.taskQueue}, BlindReviewWorkflow, BlindReviewInput{SubjectType: subjectType, SubjectID: subjectID})
	return err
}
func (s *TemporalStarter) StartReviewAudit(ctx context.Context, taskID int64) error {
	_, err := s.client.ExecuteWorkflow(ctx, client.StartWorkflowOptions{ID: fmt.Sprintf("review-audit-%d-%d", taskID, time.Now().UnixNano()), TaskQueue: s.taskQueue}, ReviewAuditWorkflow, taskID)
	return err
}
func (s *TemporalStarter) StartFinalizeReview(ctx context.Context, batchID int64) error {
	_, err := s.client.ExecuteWorkflow(ctx, client.StartWorkflowOptions{ID: fmt.Sprintf("review-fallback-%d-%d", batchID, time.Now().UnixNano()), TaskQueue: s.taskQueue}, FinalizeBlindReviewWorkflow, batchID)
	return err
}

func BlindReviewWorkflow(ctx temporalworkflow.Context, input BlindReviewInput) error {
	ctx = temporalworkflow.WithActivityOptions(ctx, temporalworkflow.ActivityOptions{StartToCloseTimeout: time.Minute, RetryPolicy: &temporal.RetryPolicy{InitialInterval: time.Second, BackoffCoefficient: 2, MaximumInterval: time.Minute, MaximumAttempts: 5}})
	var batchID int64
	if err := temporalworkflow.ExecuteActivity(ctx, "CreateBlindReviewBatch", input).Get(ctx, &batchID); err != nil {
		return err
	}
	if err := temporalworkflow.Sleep(ctx, 24*time.Hour); err != nil {
		return err
	}
	return temporalworkflow.ExecuteActivity(ctx, "FinalizeBlindReview", batchID).Get(ctx, nil)
}

func ReviewAuditWorkflow(ctx temporalworkflow.Context, taskID int64) error {
	ctx = temporalworkflow.WithActivityOptions(ctx, temporalworkflow.ActivityOptions{StartToCloseTimeout: time.Minute, RetryPolicy: &temporal.RetryPolicy{InitialInterval: 2 * time.Second, BackoffCoefficient: 2, MaximumAttempts: 5}})
	return temporalworkflow.ExecuteActivity(ctx, "AuditReview", taskID).Get(ctx, nil)
}

func FinalizeBlindReviewWorkflow(ctx temporalworkflow.Context, batchID int64) error {
	ctx = temporalworkflow.WithActivityOptions(ctx, temporalworkflow.ActivityOptions{StartToCloseTimeout: time.Minute, RetryPolicy: &temporal.RetryPolicy{InitialInterval: 2 * time.Second, BackoffCoefficient: 2, MaximumInterval: time.Minute, MaximumAttempts: 5}})
	return temporalworkflow.ExecuteActivity(ctx, "FinalizeBlindReview", batchID).Get(ctx, nil)
}

func (a *Activities) CreateBlindReviewBatch(ctx context.Context, input BlindReviewInput) (int64, error) {
	return dao.NewReviewDAO(a.DB).CreateBatch(ctx, input.SubjectType, input.SubjectID)
}

type reviewAuditResult struct {
	Fair   bool   `json:"fair"`
	Reason string `json:"reason"`
}

func (a *Activities) AuditReview(ctx context.Context, taskID int64) error {
	reviewerID, reason, result, err := dao.NewReviewDAO(a.DB).GetSubmissionForAudit(ctx, taskID)
	if err != nil {
		return err
	}
	jobKey := fmt.Sprintf("review-audit:%d", taskID)
	_, _ = a.DB.ExecContext(ctx, `INSERT INTO llm_jobs(job_key,job_type,aggregate_type,aggregate_id,status,model) VALUES($1,'review_audit','review_task',$2,'running',$3) ON CONFLICT(job_key) DO UPDATE SET status='running',attempts=llm_jobs.attempts+1`, jobKey, taskID, a.LLM.Model())
	raw, usage, err := a.LLM.CompleteJSON(ctx, `你是匿名评审公正性审核器。只返回 JSON：{"fair":boolean,"reason":string}。只评估评审是否围绕语言得体与态度真诚，不评价观点是否正确。`, fmt.Sprintf("评审结论: %s\n评审理由: %s", result, reason))
	if err != nil {
		a.recordLLMFailure(ctx, jobKey, err)
		return err
	}
	var verdict reviewAuditResult
	if err = json.Unmarshal(raw, &verdict); err != nil {
		a.recordLLMFailure(ctx, jobKey, err)
		return err
	}
	delta, status := 3, "valid"
	if !verdict.Fair {
		delta, status = -10, "invalid"
	}
	tx, err := a.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `UPDATE blind_review_tasks SET llm_check_status=$2,llm_check_result=$3 WHERE id=$1`, taskID, status, raw)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `UPDATE user_trust_profiles SET trust_score=trust_score+$2,audit_probability=LEAST(0.80,GREATEST(0.05,0.10-(trust_score+$2)/200.0)),compliant_interactions=compliant_interactions+CASE WHEN $2>0 THEN 1 ELSE 0 END,updated_at=CURRENT_TIMESTAMP WHERE user_id=$1`, reviewerID, delta)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO trust_logs(user_id,event_type,score_delta,reason,reference_type,reference_id) VALUES($1,$2,$3,$4,'review_task',$5)`, reviewerID, "review_"+status, delta, verdict.Reason, taskID)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `UPDATE llm_jobs SET status='completed',result=$2,prompt_tokens=$3,completion_tokens=$4,latency_ms=$5,completed_at=CURRENT_TIMESTAMP WHERE job_key=$1`, jobKey, raw, usage.PromptTokens, usage.CompletionTokens, usage.LatencyMS)
	if err != nil {
		return err
	}
	return tx.Commit()
}

type fallbackReviewResult struct {
	Pass   bool   `json:"pass"`
	Reason string `json:"reason"`
}

func (a *Activities) FinalizeBlindReview(ctx context.Context, batchID int64) error {
	var subjectType, status string
	var subjectID, authorID int64
	err := a.DB.QueryRowContext(ctx, `SELECT subject_type,subject_id,author_id,status FROM blind_review_batches WHERE id=$1`, batchID).Scan(&subjectType, &subjectID, &authorID, &status)
	if err != nil {
		return err
	}
	if status != "pending" {
		return nil
	}
	var subject json.RawMessage
	if subjectType == "topic" {
		err = a.DB.QueryRowContext(ctx, `SELECT jsonb_build_object('title',title,'structured_content',structured_content) FROM topics WHERE id=$1`, subjectID).Scan(&subject)
	} else {
		err = a.DB.QueryRowContext(ctx, `SELECT jsonb_build_object('statement',onboarding_statement,'background_tag',background_tag) FROM user_profiles WHERE user_id=$1`, subjectID).Scan(&subject)
	}
	if err != nil {
		return err
	}
	var reviews json.RawMessage
	err = a.DB.QueryRowContext(ctx, `SELECT COALESCE(jsonb_agg(jsonb_build_object('result',review_result,'reason',reason)) FILTER(WHERE task_status='completed'),'[]') FROM blind_review_tasks WHERE batch_id=$1`, batchID).Scan(&reviews)
	if err != nil {
		return err
	}
	jobKey := fmt.Sprintf("review-fallback:%d", batchID)
	_, _ = a.DB.ExecContext(ctx, `INSERT INTO llm_jobs(job_key,job_type,aggregate_type,aggregate_id,status,model) VALUES($1,'review_fallback','review_batch',$2,'running',$3) ON CONFLICT(job_key) DO UPDATE SET status='running',attempts=llm_jobs.attempts+1,error_message=''`, jobKey, batchID, a.LLM.Model())
	raw, usage, err := a.LLM.CompleteJSON(ctx, `你是盲审超时兜底审核器。只返回 JSON：{"pass":boolean,"reason":string}。仅判断语言是否得体、态度是否真诚。`, fmt.Sprintf("内容: %s\n已完成的人类评审: %s", subject, reviews))
	if err != nil {
		a.recordLLMFailure(ctx, jobKey, err)
		return err
	}
	var verdict fallbackReviewResult
	if err = json.Unmarshal(raw, &verdict); err != nil {
		a.recordLLMFailure(ctx, jobKey, err)
		return err
	}
	final := "reject"
	if verdict.Pass {
		final = "pass"
	}
	tx, err := a.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `UPDATE blind_review_batches SET status='expired',final_result=$2,llm_result=$3,completed_at=CURRENT_TIMESTAMP WHERE id=$1 AND status='pending'`, batchID, final, raw)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if _, err = tx.ExecContext(ctx, `UPDATE llm_jobs SET status='completed',result=$2,prompt_tokens=$3,completion_tokens=$4,latency_ms=$5,completed_at=CURRENT_TIMESTAMP,error_message='' WHERE job_key=$1`, jobKey, raw, usage.PromptTokens, usage.CompletionTokens, usage.LatencyMS); err != nil {
		return err
	}
	if affected == 0 {
		return tx.Commit()
	}
	if subjectType == "topic" {
		contentStatus := "rejected"
		if verdict.Pass {
			contentStatus = "published"
		}
		_, err = tx.ExecContext(ctx, `UPDATE topics SET status=$2,updated_at=CURRENT_TIMESTAMP WHERE id=$1 AND status='pending_review'`, subjectID, contentStatus)
	} else {
		profileStatus := "rejected"
		if verdict.Pass {
			profileStatus = "approved"
		}
		_, err = tx.ExecContext(ctx, `UPDATE user_profiles SET onboarding_status=$2,updated_at=CURRENT_TIMESTAMP WHERE user_id=$1`, authorID, profileStatus)
	}
	if err != nil {
		return err
	}
	return tx.Commit()
}
