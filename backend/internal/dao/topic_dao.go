package dao

import (
	"context"
	"database/sql"
	"errors"

	"agora-backend/internal/model"
)

type TopicDAO struct {
	db *sql.DB
}

func NewTopicDAO(db *sql.DB) *TopicDAO {
	return &TopicDAO{db: db}
}

func (d *TopicDAO) CreateTopic(ctx context.Context, t *model.Topic) error {
	query := `
		INSERT INTO topics (category_id, user_id, title, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id, view_count, post_count, like_count, created_at, updated_at
	`
	return d.db.QueryRowContext(ctx, query, t.CategoryID, t.UserID, t.Title, t.Content).
		Scan(&t.ID, &t.ViewCount, &t.PostCount, &t.LikeCount, &t.CreatedAt, &t.UpdatedAt)
}

func (d *TopicDAO) ListTopicsByCategoryID(ctx context.Context, categoryID int64, page, pageSize int) ([]*model.Topic, error) {
	// 保证无数据时返回 [] 而非 null
	topics := make([]*model.Topic, 0)

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var rows *sql.Rows
	var err error

	// 使用动态 SQL 避免 Postgres 参数类型推断问题
	if categoryID > 0 {
		query := `
			SELECT t.id, t.category_id, t.user_id, t.title, t.content, t.view_count, t.post_count, t.like_count, t.created_at, t.updated_at,
			       u.username, COALESCE(u.avatar, '')
			FROM topics t
			JOIN users u ON t.user_id = u.id
			WHERE t.category_id = $1
			ORDER BY t.created_at DESC
			LIMIT $2 OFFSET $3
		`
		rows, err = d.db.QueryContext(ctx, query, categoryID, pageSize, offset)
	} else {
		query := `
			SELECT t.id, t.category_id, t.user_id, t.title, t.content, t.view_count, t.post_count, t.like_count, t.created_at, t.updated_at,
			       u.username, COALESCE(u.avatar, '')
			FROM topics t
			JOIN users u ON t.user_id = u.id
			ORDER BY t.created_at DESC
			LIMIT $1 OFFSET $2
		`
		rows, err = d.db.QueryContext(ctx, query, pageSize, offset)
	}

	if err != nil {
		return topics, err
	}
	defer rows.Close()

	for rows.Next() {
		t := &model.Topic{}
		if err := rows.Scan(
			&t.ID, &t.CategoryID, &t.UserID, &t.Title, &t.Content, &t.ViewCount, &t.PostCount, &t.LikeCount, &t.CreatedAt, &t.UpdatedAt,
			&t.AuthorName, &t.AuthorAvatar,
		); err != nil {
			return topics, err
		}
		topics = append(topics, t)
	}

	if err := rows.Err(); err != nil {
		return topics, err
	}

	return topics, nil
}

func (d *TopicDAO) GetTopicByID(ctx context.Context, id int64) (*model.Topic, error) {
	query := `
		SELECT t.id, t.category_id, t.user_id, t.title, t.content, t.view_count, t.post_count, t.like_count, t.created_at, t.updated_at,
		       u.username, COALESCE(u.avatar, '')
		FROM topics t
		JOIN users u ON t.user_id = u.id
		WHERE t.id = $1
	`
	t := &model.Topic{}
	err := d.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.CategoryID, &t.UserID, &t.Title, &t.Content, &t.ViewCount, &t.PostCount, &t.LikeCount, &t.CreatedAt, &t.UpdatedAt,
		&t.AuthorName, &t.AuthorAvatar,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return t, err
}

func (d *TopicDAO) IncrementViewCount(ctx context.Context, id int64) error {
	query := `UPDATE topics SET view_count = view_count + 1 WHERE id = $1`
	_, err := d.db.ExecContext(ctx, query, id)
	return err
}