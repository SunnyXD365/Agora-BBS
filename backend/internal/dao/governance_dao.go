package dao

import (
	"context"
	"database/sql"
	"time"

	"agora-backend/internal/model"
)

type GovernanceDAO struct{ db *sql.DB }

func NewGovernanceDAO(db *sql.DB) *GovernanceDAO { return &GovernanceDAO{db: db} }

func (d *GovernanceDAO) StartReading(ctx context.Context, publicID string, userID, topicID int64) (*model.ReadingSession, error) {
	s := &model.ReadingSession{}
	err := d.db.QueryRowContext(ctx, `
		INSERT INTO reading_sessions(public_id, user_id, topic_id)
		SELECT $1, $2, id FROM topics WHERE id = $3 AND status = 'published'
		RETURNING public_id, topic_id, last_progress, reading_seconds, reply_dwell_seconds,
		          bottom_reached, eligible, completed, last_heartbeat_at`, publicID, userID, topicID).
		Scan(&s.PublicID, &s.TopicID, &s.Progress, &s.ReadingSeconds, &s.ReplyDwellSeconds,
			&s.BottomReached, &s.Eligible, &s.Completed, &s.LastHeartbeatAt)
	return s, err
}

func (d *GovernanceDAO) Heartbeat(ctx context.Context, publicID string, userID int64, progress int, replyFocused bool) (*model.ReadingSession, error) {
	s := &model.ReadingSession{}
	err := d.db.QueryRowContext(ctx, `
		UPDATE reading_sessions SET
			reading_seconds = reading_seconds + CASE WHEN $3 >= last_progress THEN LEAST(10, GREATEST(0, EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - last_heartbeat_at))::int)) ELSE 0 END,
			reply_dwell_seconds = reply_dwell_seconds + CASE WHEN $3 >= last_progress AND $4 THEN LEAST(10, GREATEST(0, EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - last_heartbeat_at))::int)) ELSE 0 END,
			last_progress = GREATEST(last_progress, $3),
			bottom_reached = bottom_reached OR $3 = 100,
			last_heartbeat_at = CURRENT_TIMESTAMP
		WHERE public_id = $1 AND user_id = $2 AND completed = FALSE
		RETURNING public_id, topic_id, last_progress, reading_seconds, reply_dwell_seconds,
		          bottom_reached, eligible, completed, last_heartbeat_at`, publicID, userID, progress, replyFocused).
		Scan(&s.PublicID, &s.TopicID, &s.Progress, &s.ReadingSeconds, &s.ReplyDwellSeconds,
			&s.BottomReached, &s.Eligible, &s.Completed, &s.LastHeartbeatAt)
	return s, err
}

func (d *GovernanceDAO) CompleteReading(ctx context.Context, publicID string, userID int64, dwellRequired int) (*model.ReadingSession, bool, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, false, err
	}
	defer tx.Rollback()
	s := &model.ReadingSession{}
	var wasCompleted bool
	err = tx.QueryRowContext(ctx, `
		SELECT public_id, topic_id, last_progress, reading_seconds, reply_dwell_seconds,
		       bottom_reached, eligible, completed, last_heartbeat_at
		FROM reading_sessions WHERE public_id = $1 AND user_id = $2 FOR UPDATE`, publicID, userID).
		Scan(&s.PublicID, &s.TopicID, &s.Progress, &s.ReadingSeconds, &s.ReplyDwellSeconds,
			&s.BottomReached, &s.Eligible, &wasCompleted, &s.LastHeartbeatAt)
	if err != nil {
		return nil, false, err
	}
	newlyCompleted := !wasCompleted
	if newlyCompleted {
		s.Completed = true
		s.Eligible = s.BottomReached && s.ReplyDwellSeconds >= dwellRequired
		_, err = tx.ExecContext(ctx, `UPDATE reading_sessions SET completed = TRUE, eligible = $3, completed_at = CURRENT_TIMESTAMP WHERE public_id = $1 AND user_id = $2`, publicID, userID, s.Eligible)
		if err != nil {
			return nil, false, err
		}
		_, err = tx.ExecContext(ctx, `UPDATE user_trust_profiles SET verified_read_seconds = verified_read_seconds + $2, updated_at = CURRENT_TIMESTAMP WHERE user_id = $1`, userID, s.ReadingSeconds)
		if err != nil {
			return nil, false, err
		}
	} else {
		s.Completed = true
	}
	if err := tx.Commit(); err != nil {
		return nil, false, err
	}
	return s, newlyCompleted, nil
}

func (d *GovernanceDAO) UpdateUnlockLevel(ctx context.Context, userID int64, development bool, level1, level2, level3 int) error {
	minAge1, minAge2, minAge3 := 24*time.Hour, 72*time.Hour, 168*time.Hour
	if development {
		minAge1, minAge2, minAge3 = 0, 0, 0
	}
	_, err := d.db.ExecContext(ctx, `
		UPDATE user_trust_profiles tp SET unlock_level = GREATEST(tp.unlock_level,
			CASE
				WHEN tp.verified_read_seconds >= $2 AND tp.compliant_interactions >= 5 AND tp.trust_score >= 20 AND CURRENT_TIMESTAMP - u.created_at >= $5::interval THEN 3
				WHEN tp.verified_read_seconds >= $3 AND tp.compliant_interactions >= 2 AND CURRENT_TIMESTAMP - u.created_at >= $6::interval THEN 2
				WHEN tp.verified_read_seconds >= $4 AND CURRENT_TIMESTAMP - u.created_at >= $7::interval THEN 1
				ELSE 0
			END), updated_at = CURRENT_TIMESTAMP
		FROM users u WHERE tp.user_id = u.id AND tp.user_id = $1`, userID, level3, level2, level1,
		minAge3.String(), minAge2.String(), minAge1.String())
	return err
}

func (d *GovernanceDAO) HasEligibleReading(ctx context.Context, userID, topicID int64) (bool, error) {
	var eligible bool
	err := d.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM reading_sessions WHERE user_id = $1 AND topic_id = $2 AND completed = TRUE AND eligible = TRUE)`, userID, topicID).Scan(&eligible)
	return eligible, err
}

func (d *GovernanceDAO) TopicRequiresReading(ctx context.Context, topicID int64, longChars int) (bool, error) {
	var required bool
	err := d.db.QueryRowContext(ctx, `
		SELECT c.requires_review OR char_length(COALESCE(t.structured_content->>'claim','') || COALESCE(t.structured_content->>'evidence','') || COALESCE(t.structured_content->>'uncertainty','')) >= $2
		FROM topics t JOIN categories c ON c.id = t.category_id WHERE t.id = $1 AND t.status = 'published'`, topicID, longChars).Scan(&required)
	return required, err
}

func (d *GovernanceDAO) SaveOnboarding(ctx context.Context, userID int64, statement, backgroundTag string) error {
	result, err := d.db.ExecContext(ctx, `UPDATE user_profiles SET onboarding_statement = $2, background_tag = $3, onboarding_status = 'pending_review', updated_at = CURRENT_TIMESTAMP WHERE user_id = $1`, userID, statement, backgroundTag)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
