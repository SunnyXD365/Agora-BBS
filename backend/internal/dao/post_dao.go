package dao

import (
	"context"
	"database/sql"

	"agora-backend/internal/model"
)

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

	insertQuery := `
		INSERT INTO posts (topic_id, user_id, parent_id, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id, like_count, created_at, updated_at
	`
	if err := tx.QueryRowContext(ctx, insertQuery, p.TopicID, p.UserID, p.ParentID, p.Content).
		Scan(&p.ID, &p.LikeCount, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return err
	}

	updateQuery := `UPDATE topics SET post_count = post_count + 1 WHERE id = $1`
	if _, err := tx.ExecContext(ctx, updateQuery, p.TopicID); err != nil {
		return err
	}

	return tx.Commit()
}

func (d *PostDAO) ListPostsByTopicID(ctx context.Context, topicID int64, page, pageSize int) ([]*model.Post, error) {
	posts := make([]*model.Post, 0)

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	query := `
		SELECT p.id, p.topic_id, p.user_id, p.content, p.created_at, u.username
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.topic_id = $1
		ORDER BY p.created_at ASC
		LIMIT $2 OFFSET $3
	`
	rows, err := d.db.QueryContext(ctx, query, topicID, pageSize, offset)
	if err != nil {
		return posts, err
	}
	defer rows.Close()

	for rows.Next() {
		p := &model.Post{}
		if err := rows.Scan(&p.ID, &p.TopicID, &p.UserID, &p.Content, &p.CreatedAt, &p.AuthorName); err != nil {
			return posts, err
		}
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return posts, err
	}

	return posts, nil
}