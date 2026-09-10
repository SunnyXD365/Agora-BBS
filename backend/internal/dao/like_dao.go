package dao

import (
	"context"
	"database/sql"
	"errors"

	"Agora-BBS/internal/model"
)

type LikeDAO struct {
	db *sql.DB
}

func NewLikeDAO(db *sql.DB) *LikeDAO {
	return &LikeDAO{db: db}
}

func (d *LikeDAO) GetLike(ctx context.Context, userID int64, targetType string, targetID int64) (*model.Like, error) {
	query := `
		SELECT id, user_id, target_type, target_id, created_at
		FROM likes
		WHERE user_id = $1 AND target_type = $2 AND target_id = $3
	`
	l := &model.Like{}
	err := d.db.QueryRowContext(ctx, query, userID, targetType, targetID).Scan(
		&l.ID, &l.UserID, &l.TargetType, &l.TargetID, &l.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return l, err
}

func (d *LikeDAO) CreateLike(ctx context.Context, userID int64, targetType string, targetID int64) error {
	query := `
		INSERT INTO likes (user_id, target_type, target_id)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING
	`
	_, err := d.db.ExecContext(ctx, query, userID, targetType, targetID)
	return err
}

func (d *LikeDAO) DeleteLike(ctx context.Context, userID int64, targetType string, targetID int64) error {
	query := `DELETE FROM likes WHERE user_id = $1 AND target_type = $2 AND target_id = $3`
	_, err := d.db.ExecContext(ctx, query, userID, targetType, targetID)
	return err
}
