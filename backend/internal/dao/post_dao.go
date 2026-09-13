package dao

import (
	"context"
	"database/sql"
	"errors"

	"agora-backend/internal/model"
)

var ErrInvalidParent = errors.New("parent post does not belong to topic")

type PostDAO struct {
	db *sql.DB
}

func NewPostDAO(db *sql.DB) *PostDAO {
	return &PostDAO{db: db}
}

func (d *PostDAO) CreatePost(ctx context.Context, p *model.Post) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var topicExists bool
	if err := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM topics WHERE id = $1 AND status = 'published')`, p.TopicID).Scan(&topicExists); err != nil {
		return err
	}
	if !topicExists {
		return sql.ErrNoRows
	}
	if p.ParentID != nil {
		var validParent bool
		if err := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1 AND topic_id = $2 AND status = 'published')`, *p.ParentID, p.TopicID).Scan(&validParent); err != nil {
			return err
		}
		if !validParent {
			return ErrInvalidParent
		}
	}

	insertQuery := `
		INSERT INTO posts (topic_id, user_id, parent_id, content, post_type, status, cooling_ends_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, like_count, status, cooling_ends_at, created_at, updated_at
	`
	if err := tx.QueryRowContext(ctx, insertQuery, p.TopicID, p.UserID, p.ParentID, p.Content, p.PostType, p.Status, p.CoolingEndsAt).
		Scan(&p.ID, &p.LikeCount, &p.Status, &p.CoolingEndsAt, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return err
	}

	return tx.Commit()
}

func (d *PostDAO) ListPostsByTopicID(ctx context.Context, topicID, viewerID int64, page, pageSize int) ([]*model.Post, int64, error) {
	posts := make([]*model.Post, 0)

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	query := `
		SELECT p.id, p.topic_id, p.user_id, p.parent_id, p.content, p.post_type, p.status, p.cooling_ends_at,
		       p.like_count, p.created_at, p.updated_at, u.username, COALESCE(u.avatar, '')
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.topic_id = $1 AND (p.status = 'published' OR (p.user_id = $2 AND p.status = 'cooling'))
		ORDER BY p.created_at ASC
		LIMIT $3 OFFSET $4
	`
	rows, err := d.db.QueryContext(ctx, query, topicID, viewerID, pageSize, offset)
	if err != nil {
		return posts, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		p := &model.Post{}
		if err := rows.Scan(&p.ID, &p.TopicID, &p.UserID, &p.ParentID, &p.Content, &p.PostType, &p.Status, &p.CoolingEndsAt,
			&p.LikeCount, &p.CreatedAt, &p.UpdatedAt, &p.AuthorName, &p.AuthorAvatar); err != nil {
			return posts, 0, err
		}
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return posts, 0, err
	}
	var total int64
	err = d.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM posts WHERE topic_id = $1 AND (status = 'published' OR (user_id = $2 AND status = 'cooling'))`, topicID, viewerID).Scan(&total)
	return posts, total, err
}

func (d *PostDAO) UpdateCooling(ctx context.Context, p *model.Post) error {
	result, err := d.db.ExecContext(ctx, `UPDATE posts SET content = $3, post_type = $4, cooling_ends_at = $5, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND user_id = $2 AND status = 'cooling'`, p.ID, p.UserID, p.Content, p.PostType, p.CoolingEndsAt)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (d *PostDAO) RecallCooling(ctx context.Context, id, userID int64) error {
	result, err := d.db.ExecContext(ctx, `UPDATE posts SET status = 'recalled', content = '', updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND user_id = $2 AND status = 'cooling'`, id, userID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
