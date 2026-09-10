package dao

import (
	"context"
	"database/sql"
	"errors"

	"Agora-BBS/internal/model"
)

type PostDAO struct {
	db *sql.DB
}

func NewPostDAO(db *sql.DB) *PostDAO {
	return &PostDAO{db: db}
}

func (d *PostDAO) Create(ctx context.Context, p *model.Post) error {
	query := `
		INSERT INTO posts (topic_id, user_id, parent_id, content, like_count, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	return d.db.QueryRowContext(
		ctx, query,
		p.TopicID, p.UserID, p.ParentID, p.Content, p.LikeCount, p.Status,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (d *PostDAO) ListByTopic(ctx context.Context, topicID int64, offset, limit int) ([]*model.Post, error) {
	query := `
		SELECT p.id, p.topic_id, p.user_id, p.parent_id, p.content, p.like_count, p.status,
		       p.created_at, p.updated_at, u.username, u.avatar
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.topic_id = $1 AND p.status = 'normal'
		ORDER BY p.id ASC
		LIMIT $2 OFFSET $3
	`
	rows, err := d.db.QueryContext(ctx, query, topicID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.Post
	for rows.Next() {
		p := &model.Post{}
		if err := rows.Scan(
			&p.ID, &p.TopicID, &p.UserID, &p.ParentID, &p.Content, &p.LikeCount, &p.Status,
			&p.CreatedAt, &p.UpdatedAt, &p.AuthorName, &p.AuthorAvatar,
		); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (d *PostDAO) GetByID(ctx context.Context, id int64) (*model.Post, error) {
	query := `
		SELECT id, topic_id, user_id, parent_id, content, like_count, status, created_at, updated_at
		FROM posts WHERE id = $1 AND status = 'normal'
	`
	p := &model.Post{}
	err := d.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.TopicID, &p.UserID, &p.ParentID, &p.Content, &p.LikeCount, &p.Status, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return p, err
}

func (d *PostDAO) UpdateLikeCount(ctx context.Context, postID int64, delta int) (int, error) {
	query := `
		UPDATE posts SET like_count = GREATEST(0, like_count + $1)
		WHERE id = $2
		RETURNING like_count
	`
	var count int
	err := d.db.QueryRowContext(ctx, query, delta, postID).Scan(&count)
	return count, err
}
