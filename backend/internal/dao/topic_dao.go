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
		INSERT INTO topics (category_id, user_id, title, content, structured_content, status, cooling_ends_at)
		SELECT $1, $2, $3, $4, $5, $6, $7
		WHERE EXISTS (SELECT 1 FROM categories WHERE id = $1 AND is_active = TRUE)
		RETURNING id, view_count, post_count, like_count, status, cooling_ends_at, created_at, updated_at
	`
	return d.db.QueryRowContext(ctx, query, t.CategoryID, t.UserID, t.Title, t.Content, t.StructuredContent, t.Status, t.CoolingEndsAt).
		Scan(&t.ID, &t.ViewCount, &t.PostCount, &t.LikeCount, &t.Status, &t.CoolingEndsAt, &t.CreatedAt, &t.UpdatedAt)
}

func (d *TopicDAO) CategoryRequiresReview(ctx context.Context, categoryID int64) (bool, error) {
	var requires bool
	err := d.db.QueryRowContext(ctx, `SELECT requires_review FROM categories WHERE id = $1 AND is_active = TRUE`, categoryID).Scan(&requires)
	return requires, err
}

func (d *TopicDAO) ListTopicsByCategoryID(ctx context.Context, categoryID int64, page, pageSize int) ([]*model.Topic, int64, error) {
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
			SELECT t.id, t.category_id, t.user_id, t.title, t.content, t.structured_content, t.status, t.cooling_ends_at,
			       t.view_count, t.post_count, t.like_count, t.created_at, t.updated_at,
			       u.username, COALESCE(u.avatar, '')
			FROM topics t
			JOIN users u ON t.user_id = u.id
			WHERE t.category_id = $1 AND t.status = 'published'
			ORDER BY t.created_at DESC
			LIMIT $2 OFFSET $3
		`
		rows, err = d.db.QueryContext(ctx, query, categoryID, pageSize, offset)
	} else {
		query := `
			SELECT t.id, t.category_id, t.user_id, t.title, t.content, t.structured_content, t.status, t.cooling_ends_at,
			       t.view_count, t.post_count, t.like_count, t.created_at, t.updated_at,
			       u.username, COALESCE(u.avatar, '')
			FROM topics t
			JOIN users u ON t.user_id = u.id
			WHERE t.status = 'published'
			ORDER BY t.created_at DESC
			LIMIT $1 OFFSET $2
		`
		rows, err = d.db.QueryContext(ctx, query, pageSize, offset)
	}

	if err != nil {
		return topics, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		t := &model.Topic{}
		if err := rows.Scan(
			&t.ID, &t.CategoryID, &t.UserID, &t.Title, &t.Content, &t.StructuredContent, &t.Status, &t.CoolingEndsAt,
			&t.ViewCount, &t.PostCount, &t.LikeCount, &t.CreatedAt, &t.UpdatedAt,
			&t.AuthorName, &t.AuthorAvatar,
		); err != nil {
			return topics, 0, err
		}
		topics = append(topics, t)
	}

	if err := rows.Err(); err != nil {
		return topics, 0, err
	}
	var total int64
	if categoryID > 0 {
		err = d.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM topics WHERE category_id = $1 AND status = 'published'`, categoryID).Scan(&total)
	} else {
		err = d.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM topics WHERE status = 'published'`).Scan(&total)
	}
	return topics, total, err
}

func (d *TopicDAO) GetTopicByID(ctx context.Context, id, viewerID int64) (*model.Topic, error) {
	query := `
		SELECT t.id, t.category_id, t.user_id, t.title, t.content, t.structured_content, t.status, t.cooling_ends_at,
		       t.view_count, t.post_count, t.like_count, t.created_at, t.updated_at,
		       u.username, COALESCE(u.avatar, '')
		FROM topics t
		JOIN users u ON t.user_id = u.id
		WHERE t.id = $1 AND (t.status = 'published' OR (t.user_id = $2 AND t.status IN ('cooling', 'pending_review', 'rejected')))
	`
	t := &model.Topic{}
	err := d.db.QueryRowContext(ctx, query, id, viewerID).Scan(
		&t.ID, &t.CategoryID, &t.UserID, &t.Title, &t.Content, &t.StructuredContent, &t.Status, &t.CoolingEndsAt,
		&t.ViewCount, &t.PostCount, &t.LikeCount, &t.CreatedAt, &t.UpdatedAt,
		&t.AuthorName, &t.AuthorAvatar,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return t, err
}

func (d *TopicDAO) UpdateCooling(ctx context.Context, t *model.Topic) error {
	result, err := d.db.ExecContext(ctx, `UPDATE topics SET title = $3, content = $4, structured_content = $5, cooling_ends_at = $6, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND user_id = $2 AND status = 'cooling'`, t.ID, t.UserID, t.Title, t.Content, t.StructuredContent, t.CoolingEndsAt)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (d *TopicDAO) RecallCooling(ctx context.Context, id, userID int64) error {
	result, err := d.db.ExecContext(ctx, `UPDATE topics SET status = 'recalled', title = '[已撤回]', content = '', structured_content = '{}'::jsonb, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND user_id = $2 AND status = 'cooling'`, id, userID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (d *TopicDAO) IncrementViewCount(ctx context.Context, id int64) error {
	query := `UPDATE topics SET view_count = view_count + 1 WHERE id = $1`
	_, err := d.db.ExecContext(ctx, query, id)
	return err
}
