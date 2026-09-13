package dao

import (
	"context"
	"database/sql"

	"agora-backend/internal/model"
)

type BookmarkDAO struct{ db *sql.DB }

func NewBookmarkDAO(db *sql.DB) *BookmarkDAO { return &BookmarkDAO{db: db} }

func (d *BookmarkDAO) Create(ctx context.Context, userID, topicID int64) error {
	result, err := d.db.ExecContext(ctx, `
		INSERT INTO bookmarks (user_id, topic_id)
		SELECT $1, id FROM topics WHERE id = $2 AND status = 'published'
		ON CONFLICT (user_id, topic_id) DO NOTHING`, userID, topicID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err == nil && rows == 0 {
		var exists bool
		if scanErr := d.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM topics WHERE id = $1 AND status = 'published')`, topicID).Scan(&exists); scanErr != nil {
			return scanErr
		}
		if !exists {
			return sql.ErrNoRows
		}
	}
	return nil
}

func (d *BookmarkDAO) Delete(ctx context.Context, userID, topicID int64) error {
	_, err := d.db.ExecContext(ctx, `DELETE FROM bookmarks WHERE user_id = $1 AND topic_id = $2`, userID, topicID)
	return err
}

func (d *BookmarkDAO) List(ctx context.Context, userID int64, page, pageSize int) ([]*model.Bookmark, int64, error) {
	offset := (page - 1) * pageSize
	rows, err := d.db.QueryContext(ctx, `
		SELECT b.id, b.user_id, b.topic_id, b.created_at,
		       t.id, t.category_id, t.user_id, t.title, t.content, t.structured_content, t.status, t.cooling_ends_at,
		       t.view_count, t.post_count, t.like_count, t.created_at, t.updated_at, u.username, COALESCE(u.avatar, '')
		FROM bookmarks b
		JOIN topics t ON t.id = b.topic_id AND t.status = 'published'
		JOIN users u ON u.id = t.user_id
		WHERE b.user_id = $1 ORDER BY b.created_at DESC LIMIT $2 OFFSET $3`, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]*model.Bookmark, 0)
	for rows.Next() {
		b := &model.Bookmark{Topic: &model.Topic{}}
		if err := rows.Scan(&b.ID, &b.UserID, &b.TopicID, &b.CreatedAt,
			&b.Topic.ID, &b.Topic.CategoryID, &b.Topic.UserID, &b.Topic.Title, &b.Topic.Content,
			&b.Topic.StructuredContent, &b.Topic.Status, &b.Topic.CoolingEndsAt, &b.Topic.ViewCount,
			&b.Topic.PostCount, &b.Topic.LikeCount, &b.Topic.CreatedAt, &b.Topic.UpdatedAt,
			&b.Topic.AuthorName, &b.Topic.AuthorAvatar); err != nil {
			return nil, 0, err
		}
		items = append(items, b)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	var total int64
	err = d.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM bookmarks WHERE user_id = $1`, userID).Scan(&total)
	return items, total, err
}
