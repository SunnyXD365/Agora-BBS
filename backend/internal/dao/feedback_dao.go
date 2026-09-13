package dao

import (
	"context"
	"database/sql"
	"errors"

	"agora-backend/internal/model"
	"github.com/lib/pq"
)

var ErrSelfFeedback = errors.New("cannot give feedback to your own content")

type FeedbackDAO struct{ db *sql.DB }

func NewFeedbackDAO(db *sql.DB) *FeedbackDAO { return &FeedbackDAO{db: db} }

func (d *FeedbackDAO) Upsert(ctx context.Context, userID int64, req *model.UpsertFeedbackReq) (*model.Feedback, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var authorID int64
	if req.TargetType == "topic" {
		err = tx.QueryRowContext(ctx, `SELECT user_id FROM topics WHERE id = $1 AND status = 'published'`, req.TargetID).Scan(&authorID)
	} else {
		err = tx.QueryRowContext(ctx, `SELECT user_id FROM posts WHERE id = $1 AND status = 'published'`, req.TargetID).Scan(&authorID)
	}
	if err != nil {
		return nil, err
	}
	if authorID == userID {
		return nil, ErrSelfFeedback
	}
	f := &model.Feedback{}
	err = tx.QueryRowContext(ctx, `
		INSERT INTO contextual_feedbacks(user_id, target_type, target_id, stance, tag, reason)
		VALUES($1,$2,$3,$4,$5,$6)
		ON CONFLICT(user_id,target_type,target_id) DO UPDATE SET stance=EXCLUDED.stance, tag=EXCLUDED.tag,
			reason=EXCLUDED.reason, status='published', llm_audit_status='not_selected', llm_audit_result='{}'::jsonb, updated_at=CURRENT_TIMESTAMP
		RETURNING id,user_id,target_type,target_id,stance,tag,reason,status,llm_audit_status,llm_audit_result,created_at,updated_at`,
		userID, req.TargetType, req.TargetID, req.Stance, req.Tag, req.Reason).
		Scan(&f.ID, &f.UserID, &f.TargetType, &f.TargetID, &f.Stance, &f.Tag, &f.Reason, &f.Status, &f.LLMAuditStatus, &f.LLMAuditResult, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if err := recomputeFeedbackScore(ctx, tx, req.TargetType, req.TargetID); err != nil {
		return nil, err
	}
	return f, tx.Commit()
}

func (d *FeedbackDAO) Withdraw(ctx context.Context, userID, id int64) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var targetType string
	var targetID int64
	err = tx.QueryRowContext(ctx, `UPDATE contextual_feedbacks SET status='withdrawn', updated_at=CURRENT_TIMESTAMP WHERE id=$1 AND user_id=$2 AND status='published' RETURNING target_type,target_id`, id, userID).Scan(&targetType, &targetID)
	if err != nil {
		return err
	}
	if err := recomputeFeedbackScore(ctx, tx, targetType, targetID); err != nil {
		return err
	}
	return tx.Commit()
}

func recomputeFeedbackScore(ctx context.Context, tx *sql.Tx, targetType string, targetID int64) error {
	query := `UPDATE topics SET feedback_score=(SELECT COALESCE(SUM(CASE stance WHEN 'support' THEN 1 ELSE -1 END),0) FROM contextual_feedbacks WHERE target_type='topic' AND target_id=$1 AND status='published') WHERE id=$1`
	if targetType == "post" {
		query = `UPDATE posts SET feedback_score=(SELECT COALESCE(SUM(CASE stance WHEN 'support' THEN 1 ELSE -1 END),0) FROM contextual_feedbacks WHERE target_type='post' AND target_id=$1 AND status='published') WHERE id=$1`
	}
	_, err := tx.ExecContext(ctx, query, targetID)
	return err
}

func (d *FeedbackDAO) Summary(ctx context.Context, targetType string, targetID, viewerID int64) (*model.FeedbackSummary, error) {
	summary := &model.FeedbackSummary{Tags: map[string]int{}}
	rows, err := d.db.QueryContext(ctx, `SELECT stance,tag,COUNT(*) FROM contextual_feedbacks WHERE target_type=$1 AND target_id=$2 AND status='published' GROUP BY stance,tag`, targetType, targetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var stance, tag string
		var count int
		if err := rows.Scan(&stance, &tag, &count); err != nil {
			return nil, err
		}
		summary.Tags[tag] += count
		if stance == "support" {
			summary.Support += count
		} else {
			summary.Challenge += count
		}
	}
	summary.Score = summary.Support - summary.Challenge
	if viewerID > 0 {
		f := &model.Feedback{}
		err := d.db.QueryRowContext(ctx, `SELECT id,user_id,target_type,target_id,stance,tag,reason,status,llm_audit_status,llm_audit_result,created_at,updated_at FROM contextual_feedbacks WHERE user_id=$1 AND target_type=$2 AND target_id=$3 AND status='published'`, viewerID, targetType, targetID).Scan(&f.ID, &f.UserID, &f.TargetType, &f.TargetID, &f.Stance, &f.Tag, &f.Reason, &f.Status, &f.LLMAuditStatus, &f.LLMAuditResult, &f.CreatedAt, &f.UpdatedAt)
		if err == nil {
			summary.Mine = f
		} else if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	return summary, rows.Err()
}

func (d *FeedbackDAO) AuditProbability(ctx context.Context, userID int64) (float64, error) {
	var p float64
	err := d.db.QueryRowContext(ctx, `SELECT audit_probability FROM user_trust_profiles WHERE user_id=$1`, userID).Scan(&p)
	return p, err
}

func (d *FeedbackDAO) QueueAudit(ctx context.Context, id int64) error {
	_, err := d.db.ExecContext(ctx, `UPDATE contextual_feedbacks SET llm_audit_status='pending' WHERE id=$1`, id)
	return err
}

func (d *FeedbackDAO) ListClusters(ctx context.Context, topicID int64) ([]*model.CommentCluster, error) {
	rows, err := d.db.QueryContext(ctx, `SELECT c.id,c.topic_id,c.tag,c.summary,c.weight,c.generation,c.created_at,COALESCE(array_agg(a.post_id) FILTER (WHERE a.post_id IS NOT NULL),'{}') FROM comment_clusters c LEFT JOIN post_cluster_assignments a ON a.cluster_id=c.id WHERE c.topic_id=$1 GROUP BY c.id ORDER BY c.weight DESC,c.id`, topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]*model.CommentCluster, 0)
	for rows.Next() {
		c := &model.CommentCluster{}
		if err := rows.Scan(&c.ID, &c.TopicID, &c.Tag, &c.Summary, &c.Weight, &c.Generation, &c.CreatedAt, pq.Array(&c.PostIDs)); err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	return items, rows.Err()
}
