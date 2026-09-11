package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	temporalworkflow "go.temporal.io/sdk/workflow"
)

type LLMStarter interface {
	StartFeedbackAudit(context.Context, int64) error
	StartClusterTopic(context.Context, int64) error
}

func (s *TemporalStarter) StartFeedbackAudit(ctx context.Context, feedbackID int64) error {
	_, err := s.client.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		ID: fmt.Sprintf("feedback-audit-%d-%d", feedbackID, time.Now().UnixNano()), TaskQueue: s.taskQueue,
	}, FeedbackAuditWorkflow, feedbackID)
	return err
}

func (s *TemporalStarter) StartClusterTopic(ctx context.Context, topicID int64) error {
	_, err := s.client.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		ID: fmt.Sprintf("comment-cluster-%d-%d", topicID, time.Now().UnixNano()), TaskQueue: s.taskQueue,
	}, ClusterTopicWorkflow, topicID)
	return err
}

func ClusterTopicWorkflow(ctx temporalworkflow.Context, topicID int64) error {
	ctx = temporalworkflow.WithActivityOptions(ctx, temporalworkflow.ActivityOptions{
		StartToCloseTimeout: 2 * time.Minute,
		RetryPolicy:         &temporal.RetryPolicy{InitialInterval: 2 * time.Second, BackoffCoefficient: 2, MaximumInterval: time.Minute, MaximumAttempts: 5},
	})
	return temporalworkflow.ExecuteActivity(ctx, "ClusterTopic", topicID).Get(ctx, nil)
}

func FeedbackAuditWorkflow(ctx temporalworkflow.Context, feedbackID int64) error {
	ctx = temporalworkflow.WithActivityOptions(ctx, temporalworkflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy:         &temporal.RetryPolicy{InitialInterval: 2 * time.Second, BackoffCoefficient: 2, MaximumInterval: time.Minute, MaximumAttempts: 5},
	})
	return temporalworkflow.ExecuteActivity(ctx, "AuditFeedback", feedbackID).Get(ctx, nil)
}

type feedbackAuditResult struct {
	Valid  bool   `json:"valid"`
	Reason string `json:"reason"`
}

func (a *Activities) AuditFeedback(ctx context.Context, feedbackID int64) error {
	var userID, targetID int64
	var targetType, stance, tag, reason string
	err := a.DB.QueryRowContext(ctx, `SELECT user_id,target_type,target_id,stance,tag,reason FROM contextual_feedbacks WHERE id=$1 AND status='published'`, feedbackID).
		Scan(&userID, &targetType, &targetID, &stance, &tag, &reason)
	if err != nil {
		return err
	}
	jobKey := fmt.Sprintf("feedback-audit:%d", feedbackID)
	_, _ = a.DB.ExecContext(ctx, `INSERT INTO llm_jobs(job_key,job_type,aggregate_type,aggregate_id,status,model) VALUES($1,'feedback_audit','feedback',$2,'running',$3) ON CONFLICT(job_key) DO UPDATE SET status='running',attempts=llm_jobs.attempts+1,error_message=''`, jobKey, feedbackID, a.LLM.Model())
	prompt := fmt.Sprintf("立场: %s\n标签: %s\n理由: %s", stance, tag, reason)
	raw, usage, err := a.LLM.CompleteJSON(ctx,
		`你是社区互动反作弊审核器。只返回 JSON：{"valid":boolean,"reason":string}。判断理由是否具体、与标签相关、非垃圾文本或恶意群体攻击。`, prompt)
	if err != nil {
		a.recordLLMFailure(ctx, jobKey, err)
		return err
	}
	var result feedbackAuditResult
	if err := json.Unmarshal(raw, &result); err != nil {
		a.recordLLMFailure(ctx, jobKey, err)
		return err
	}
	tx, err := a.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	delta, auditStatus := 1, "valid"
	if !result.Valid {
		delta, auditStatus = -5, "invalid"
	}
	_, err = tx.ExecContext(ctx, `UPDATE contextual_feedbacks SET llm_audit_status=$2,llm_audit_result=$3,status=CASE WHEN $4 THEN status ELSE 'rejected' END,updated_at=CURRENT_TIMESTAMP WHERE id=$1`, feedbackID, auditStatus, raw, result.Valid)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `UPDATE user_trust_profiles SET trust_score=trust_score+$2,audit_probability=LEAST(0.80,GREATEST(0.05,0.10-(trust_score+$2)/200.0)),compliant_interactions=compliant_interactions+CASE WHEN $2>0 THEN 1 ELSE 0 END,updated_at=CURRENT_TIMESTAMP WHERE user_id=$1`, userID, delta)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO trust_logs(user_id,event_type,score_delta,reason,reference_type,reference_id) VALUES($1,$2,$3,$4,'feedback',$5)`, userID, "feedback_"+auditStatus, delta, result.Reason, feedbackID)
	if err != nil {
		return err
	}
	if !result.Valid {
		query := `UPDATE topics SET feedback_score=(SELECT COALESCE(SUM(CASE stance WHEN 'support' THEN 1 ELSE -1 END),0) FROM contextual_feedbacks WHERE target_type='topic' AND target_id=$1 AND status='published') WHERE id=$1`
		if targetType == "post" {
			query = `UPDATE posts SET feedback_score=(SELECT COALESCE(SUM(CASE stance WHEN 'support' THEN 1 ELSE -1 END),0) FROM contextual_feedbacks WHERE target_type='post' AND target_id=$1 AND status='published') WHERE id=$1`
		}
		if _, err = tx.ExecContext(ctx, query, targetID); err != nil {
			return err
		}
	}
	_, err = tx.ExecContext(ctx, `UPDATE llm_jobs SET status='completed',result=$2,prompt_tokens=$3,completion_tokens=$4,latency_ms=$5,completed_at=CURRENT_TIMESTAMP WHERE job_key=$1`, jobKey, raw, usage.PromptTokens, usage.CompletionTokens, usage.LatencyMS)
	if err != nil {
		return err
	}
	return tx.Commit()
}

type clusterResponse struct {
	Clusters []struct {
		Tag     string  `json:"tag"`
		Summary string  `json:"summary"`
		PostIDs []int64 `json:"post_ids"`
		Weight  float64 `json:"weight"`
	} `json:"clusters"`
}

func (a *Activities) ClusterTopic(ctx context.Context, topicID int64) error {
	rows, err := a.DB.QueryContext(ctx, `SELECT id,post_type,content,feedback_score FROM posts WHERE topic_id=$1 AND status='published' ORDER BY created_at`, topicID)
	if err != nil {
		return err
	}
	defer rows.Close()
	type post struct {
		ID      int64  `json:"id"`
		Type    string `json:"type"`
		Content string `json:"content"`
		Score   int    `json:"score"`
	}
	posts := make([]post, 0)
	for rows.Next() {
		var p post
		if err := rows.Scan(&p.ID, &p.Type, &p.Content, &p.Score); err != nil {
			return err
		}
		posts = append(posts, p)
	}
	if len(posts) < 20 {
		return nil
	}
	payload, _ := json.Marshal(posts)
	jobKey := fmt.Sprintf("comment-cluster:%d:%d", topicID, len(posts))
	_, _ = a.DB.ExecContext(ctx, `INSERT INTO llm_jobs(job_key,job_type,aggregate_type,aggregate_id,status,model) VALUES($1,'comment_cluster','topic',$2,'running',$3) ON CONFLICT(job_key) DO UPDATE SET status='running',attempts=llm_jobs.attempts+1`, jobKey, topicID, a.LLM.Model())
	raw, usage, err := a.LLM.CompleteJSON(ctx,
		`将讨论归纳为 3 到 8 个互不重复的中文标签。只返回 JSON：{"clusters":[{"tag":string,"summary":string,"post_ids":[number],"weight":number}]}。post_ids 只能使用输入 ID。`, string(payload))
	if err != nil {
		a.recordLLMFailure(ctx, jobKey, err)
		return err
	}
	var result clusterResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		a.recordLLMFailure(ctx, jobKey, err)
		return err
	}
	if len(result.Clusters) < 3 || len(result.Clusters) > 8 {
		err = errors.New("LLM cluster count must be between 3 and 8")
		a.recordLLMFailure(ctx, jobKey, err)
		return err
	}
	valid := map[int64]bool{}
	for _, p := range posts {
		valid[p.ID] = true
	}
	tx, err := a.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var generation int
	_ = tx.QueryRowContext(ctx, `SELECT COALESCE(MAX(generation),0)+1 FROM comment_clusters WHERE topic_id=$1`, topicID).Scan(&generation)
	_, _ = tx.ExecContext(ctx, `DELETE FROM comment_clusters WHERE topic_id=$1`, topicID)
	for _, cluster := range result.Clusters {
		cluster.Tag = strings.TrimSpace(cluster.Tag)
		if cluster.Tag == "" {
			err = errors.New("empty cluster tag")
			a.recordLLMFailure(ctx, jobKey, err)
			return err
		}
		var clusterID int64
		err = tx.QueryRowContext(ctx, `INSERT INTO comment_clusters(topic_id,tag,summary,weight,model,generation) VALUES($1,$2,$3,$4,$5,$6) RETURNING id`, topicID, cluster.Tag, cluster.Summary, cluster.Weight, a.LLM.Model(), generation).Scan(&clusterID)
		if err != nil {
			return err
		}
		for _, postID := range cluster.PostIDs {
			if valid[postID] {
				if _, err = tx.ExecContext(ctx, `INSERT INTO post_cluster_assignments(cluster_id,post_id,relevance) VALUES($1,$2,1) ON CONFLICT DO NOTHING`, clusterID, postID); err != nil {
					return err
				}
			}
		}
	}
	_, err = tx.ExecContext(ctx, `UPDATE llm_jobs SET status='completed',result=$2,prompt_tokens=$3,completion_tokens=$4,latency_ms=$5,completed_at=CURRENT_TIMESTAMP WHERE job_key=$1`, jobKey, raw, usage.PromptTokens, usage.CompletionTokens, usage.LatencyMS)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (a *Activities) recordLLMFailure(ctx context.Context, jobKey string, jobErr error) {
	message := jobErr.Error()
	if len(message) > 2000 {
		message = message[:2000]
	}
	_, _ = a.DB.ExecContext(ctx, `UPDATE llm_jobs SET status='failed',error_message=$2 WHERE job_key=$1`, jobKey, message)
}
