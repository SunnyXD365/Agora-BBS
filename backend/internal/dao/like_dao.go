package dao

import (
	"context"
	"database/sql"

	"agora-backend/internal/model"
)

type LikeDAO struct {
	db *sql.DB
}

func NewLikeDAO(db *sql.DB) *LikeDAO {
	return &LikeDAO{db: db}
}

func (d *LikeDAO) CreateLike(ctx context.Context, like *model.Like) error {
	query := `INSERT INTO likes (user_id, target_type, target_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	_, err := d.db.ExecContext(ctx, query, like.UserID, like.TargetType, like.TargetID)
	return err
}

func (d *LikeDAO) DeleteLike(ctx context.Context, like *model.Like) error {
	query := `DELETE FROM likes WHERE user_id = $1 AND target_type = $2 AND target_id = $3`
	_, err := d.db.ExecContext(ctx, query, like.UserID, like.TargetType, like.TargetID)
	return err
}