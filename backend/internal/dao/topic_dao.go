package dao

import (
	"context"
	"database/sql"
	"errors"

	"Agora-BBS/internal/model"
)

type TopicDAO struct {
	db *sql.DB
}

func NewTopicDAO(db *sql.DB) *TopicDAO {
	return &TopicDAO{db: db}
}

func (d *TopicDAO) Create(ctx context.Context, t *model.Topic) error {
	query := `
		INSERT INTO topics (category_id, user_id, title, content, view_count, post_count, is_sticky, is_essence, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`
	return d.db.QueryRowContext(
		ctx, query,
		t.CategoryID, t.UserID, t.Title, t.Content, t.ViewCount, t.PostCount, t.IsSticky, t.IsEssence, t.Status,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (d *TopicDAO) List(ctx context.Context, categoryID int64, offset, limit int) ([]*model.Topic, error) {
	var rows *sql.Rows
	var err error

	if categoryID > 0 {
		query := `
			SELECT t.id, t.category_id, t.user_id, t.title, t.content, t.view_count, t.post_count,
			       t.is_sticky, t.is_essence, t.status, t.created_at, t.updated_at, u.username, u.avatar
			FROM topics t
			JOIN users u ON t.user_id = u.id
			WHERE t.category_id = $1 AND t.status = 'normal'
			ORDER BY t.is_sticky DESC, t.id DESC
			LIMIT $2 OFFSET $3
		`
		rows, err = d.db.QueryContext(ctx, query, categoryID, limit, offset)
	} else {
		query := `
			SELECT t.id, t.category_id, t.user_id, t.title, t.content, t.view_count, t.post_count,
			       t.is_sticky, t.is_essence, t.status, t.created_at, t.updated_at, u.username, u.avatar
			FROM topics t
			JOIN users u ON t.user_id = u.id
			WHERE t.status = 'normal'
			ORDER BY t.is_sticky DESC, t.id DESC
			LIMIT $1 OFFSET $2
		`
		rows, err = d.db.QueryContext(ctx, query, limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.Topic
	for rows.Next() {
		t := &model.Topic{}
		if err := rows.Scan(
			&t.ID, &t.CategoryID, &t.UserID, &t.Title, &t.Content, &t.ViewCount, &t.PostCount,
			&t.IsSticky, &t.IsEssence, &t.Status, &t.CreatedAt, &t.UpdatedAt, &t.AuthorName, &t.AuthorAvatar,
		); err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, nil
}

func (d *TopicDAO) GetByID(ctx context.Context, id int64) (*model.Topic, error) {
	query := `
		SELECT t.id, t.category_id, t.user_id, t.title, t.content, t.view_count, t.post_count,
		       t.is_sticky, t.is_essence, t.status, t.created_at, t.updated_at, u.username, u.avatar
		FROM topics t
		JOIN users u ON t.user_id = u.id
		WHERE t.id = $1 AND t.status = 'normal'
	`
	t := &model.Topic{}
	err := d.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.CategoryID, &t.UserID, &t.Title, &t.Content, &t.ViewCount, &t.PostCount,
		&t.IsSticky, &t.IsEssence, &t.Status, &t.CreatedAt, &t.UpdatedAt, &t.AuthorName, &t.AuthorAvatar,
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

func (d *TopicDAO) IncrementPostCount(ctx context.Context, id int64) error {
	query := `UPDATE topics SET post_count = post_count + 1 WHERE id = $1`
	_, err := d.db.ExecContext(ctx, query, id)
	return err
}

func (d *TopicDAO) UpdateLikeCount(ctx context.Context, topicID int64, delta int) (int, error) {
	query := `
		UPDATE topics SET like_count = GREATEST(0, like_count + $1)
		WHERE id = $2
		RETURNING like_count
	`
	var count int
	err := d.db.QueryRowContext(ctx, query, delta, topicID).Scan(&count)
	return count, err
}
