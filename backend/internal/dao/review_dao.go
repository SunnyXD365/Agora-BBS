package dao

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"agora-backend/internal/model"
)

type ReviewDAO struct{ db *sql.DB }

func NewReviewDAO(db *sql.DB) *ReviewDAO { return &ReviewDAO{db: db} }

func (d *ReviewDAO) ListTasks(ctx context.Context, reviewerID int64) ([]*model.ReviewTask, error) {
	rows, err := d.db.QueryContext(ctx, `
		SELECT t.id,t.batch_id,b.subject_type,
		CASE WHEN b.subject_type='user' THEN jsonb_build_object('statement',p.onboarding_statement,'background_tag',p.background_tag)
		     ELSE jsonb_build_object('title',tp.title,'claim',tp.structured_content->>'claim','evidence',tp.structured_content->>'evidence','uncertainty',tp.structured_content->>'uncertainty') END,
		t.task_status,b.deadline,t.created_at
		FROM blind_review_tasks t JOIN blind_review_batches b ON b.id=t.batch_id
		LEFT JOIN user_profiles p ON b.subject_type='user' AND p.user_id=b.subject_id
		LEFT JOIN topics tp ON b.subject_type='topic' AND tp.id=b.subject_id
		WHERE t.reviewer_id=$1 AND t.task_status='pending' AND b.status='pending' ORDER BY b.deadline`, reviewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]*model.ReviewTask, 0)
	for rows.Next() {
		task := &model.ReviewTask{}
		if err := rows.Scan(&task.ID, &task.BatchID, &task.SubjectType, &task.Subject, &task.TaskStatus, &task.Deadline, &task.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, task)
	}
	return items, rows.Err()
}

func (d *ReviewDAO) Submit(ctx context.Context, reviewerID, taskID int64, req *model.SubmitReviewReq) (*model.ReviewSubmission, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	result := "reject"
	if req.Appropriateness && req.Sincerity {
		result = "pass"
	}
	s := &model.ReviewSubmission{}
	err = tx.QueryRowContext(ctx, `UPDATE blind_review_tasks SET task_status='completed',appropriateness=$3,sincerity=$4,review_result=$5,reason=$6,completed_at=CURRENT_TIMESTAMP WHERE id=$1 AND reviewer_id=$2 AND task_status='pending' RETURNING id,batch_id,review_result,llm_check_status,completed_at`, taskID, reviewerID, req.Appropriateness, req.Sincerity, result, req.Reason).Scan(&s.ID, &s.BatchID, &s.Result, &s.LLMCheckStatus, &s.CompletedAt)
	if err != nil {
		return nil, err
	}
	var passCount, rejectCount int
	err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FILTER(WHERE review_result='pass'),COUNT(*) FILTER(WHERE review_result='reject') FROM blind_review_tasks WHERE batch_id=$1 AND task_status='completed'`, s.BatchID).Scan(&passCount, &rejectCount)
	if err != nil {
		return nil, err
	}
	final := ""
	if passCount >= 2 {
		final = "pass"
	} else if rejectCount >= 2 {
		final = "reject"
	}
	if final != "" {
		var subjectType string
		var subjectID, authorID int64
		err = tx.QueryRowContext(ctx, `UPDATE blind_review_batches SET status='completed',final_result=$2,completed_at=CURRENT_TIMESTAMP WHERE id=$1 AND status='pending' RETURNING subject_type,subject_id,author_id`, s.BatchID, final).Scan(&subjectType, &subjectID, &authorID)
		if err == nil {
			if subjectType == "topic" {
				status := "rejected"
				if final == "pass" {
					status = "published"
				}
				_, err = tx.ExecContext(ctx, `UPDATE topics SET status=$2,updated_at=CURRENT_TIMESTAMP WHERE id=$1 AND status='pending_review'`, subjectID, status)
			} else {
				status := "rejected"
				delta := -5
				if final == "pass" {
					status = "approved"
					delta = 5
				}
				_, err = tx.ExecContext(ctx, `UPDATE user_profiles SET onboarding_status=$2,updated_at=CURRENT_TIMESTAMP WHERE user_id=$1`, authorID, status)
				if err == nil {
					_, err = tx.ExecContext(ctx, `UPDATE user_trust_profiles SET trust_score=trust_score+$2,audit_probability=LEAST(0.80,GREATEST(0.05,0.10-(trust_score+$2)/200.0)),updated_at=CURRENT_TIMESTAMP WHERE user_id=$1`, authorID, delta)
				}
				if err == nil {
					_, err = tx.ExecContext(ctx, `INSERT INTO trust_logs(user_id,event_type,score_delta,reason,reference_type,reference_id) VALUES($1,$2,$3,$4,'review_batch',$5)`, authorID, "onboarding_"+status, delta, "社区自述匿名评审结论："+final, s.BatchID)
				}
			}
		}
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	return s, tx.Commit()
}

func (d *ReviewDAO) CreateBatch(ctx context.Context, subjectType string, subjectID int64) (int64, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	var authorID int64
	if subjectType == "user" {
		authorID = subjectID
	} else {
		if err = tx.QueryRowContext(ctx, `SELECT user_id FROM topics WHERE id=$1`, subjectID).Scan(&authorID); err != nil {
			return 0, err
		}
	}
	var batchID int64
	deadline := time.Now().Add(24 * time.Hour)
	err = tx.QueryRowContext(ctx, `INSERT INTO blind_review_batches(subject_type,subject_id,author_id,deadline) VALUES($1,$2,$3,$4) ON CONFLICT(subject_type,subject_id) DO UPDATE SET deadline=EXCLUDED.deadline RETURNING id`, subjectType, subjectID, authorID, deadline).Scan(&batchID)
	if err != nil {
		return 0, err
	}
	rows, err := tx.QueryContext(ctx, `WITH eligible AS (SELECT u.id,ROW_NUMBER() OVER(PARTITION BY COALESCE(NULLIF(p.background_tag,''),u.id::text) ORDER BY random()) AS rn FROM users u JOIN user_trust_profiles tp ON tp.user_id=u.id LEFT JOIN user_profiles p ON p.user_id=u.id WHERE u.status='active' AND tp.unlock_level>=3 AND u.id<>$1) SELECT id FROM eligible ORDER BY rn,random() LIMIT 3`, authorID)
	if err != nil {
		return 0, err
	}
	reviewers := make([]int64, 0, 3)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return 0, err
		}
		reviewers = append(reviewers, id)
	}
	rows.Close()
	for _, reviewerID := range reviewers {
		if _, err = tx.ExecContext(ctx, `INSERT INTO blind_review_tasks(batch_id,reviewer_id) VALUES($1,$2) ON CONFLICT DO NOTHING`, batchID, reviewerID); err != nil {
			return 0, err
		}
	}
	return batchID, tx.Commit()
}

func (d *ReviewDAO) GetSubmissionForAudit(ctx context.Context, taskID int64) (reviewerID int64, reason, result string, err error) {
	err = d.db.QueryRowContext(ctx, `SELECT reviewer_id,reason,review_result FROM blind_review_tasks WHERE id=$1 AND task_status='completed'`, taskID).Scan(&reviewerID, &reason, &result)
	return
}
