package dao

import (
	"context"
	"database/sql"

	"agora-backend/internal/model"
)

type ContentDAO struct {
	db *sql.DB
}

func NewContentDAO(db *sql.DB) *ContentDAO {
	return &ContentDAO{db: db}
}

func (d *ContentDAO) ListMine(ctx context.Context, userID int64, contentType, status string, page, pageSize int) ([]*model.MyContentItem, int64, error) {
	items := make([]*model.MyContentItem, 0)
	offset := (page - 1) * pageSize
	var rows *sql.Rows
	var err error
	if contentType == "post" {
		rows, err = d.db.QueryContext(ctx, `
			SELECT p.id, 'post', p.topic_id, t.title, LEFT(p.content, 400), p.status, p.post_type,
			       p.cooling_ends_at, p.created_at, p.updated_at
			FROM posts p
			JOIN topics t ON t.id = p.topic_id
			WHERE p.user_id = $1 AND ($2 = '' OR p.status = $2)
			ORDER BY p.updated_at DESC
			LIMIT $3 OFFSET $4`, userID, status, pageSize, offset)
	} else {
		rows, err = d.db.QueryContext(ctx, `
			SELECT t.id, 'topic', t.id, t.title,
			       LEFT(COALESCE(NULLIF(t.structured_content->>'claim', ''), t.content), 400),
			       t.status, '', t.cooling_ends_at, t.created_at, t.updated_at
			FROM topics t
			WHERE t.user_id = $1 AND ($2 = '' OR t.status = $2)
			ORDER BY t.updated_at DESC
			LIMIT $3 OFFSET $4`, userID, status, pageSize, offset)
	}
	if err != nil {
		return items, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		item := &model.MyContentItem{}
		if err := rows.Scan(&item.ID, &item.Type, &item.TopicID, &item.Title, &item.Excerpt, &item.Status, &item.PostType, &item.CoolingEndsAt, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return items, 0, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return items, 0, err
	}
	var total int64
	if contentType == "post" {
		err = d.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM posts WHERE user_id = $1 AND ($2 = '' OR status = $2)`, userID, status).Scan(&total)
	} else {
		err = d.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM topics WHERE user_id = $1 AND ($2 = '' OR status = $2)`, userID, status).Scan(&total)
	}
	return items, total, err
}

func (d *ContentDAO) Search(ctx context.Context, pattern, contentType string, page, pageSize int) ([]*model.SearchResult, int64, error) {
	items := make([]*model.SearchResult, 0)
	offset := (page - 1) * pageSize
	typeFilter := contentType
	query := `
		WITH matched AS (
			SELECT t.id, 'topic'::text AS type, t.id AS topic_id, t.title,
			       LEFT(COALESCE(NULLIF(t.structured_content->>'claim', ''), t.content), 400) AS excerpt,
			       u.username AS author_name, t.created_at,
			       CASE WHEN t.title ILIKE $1 THEN 0 ELSE 1 END AS rank
			FROM topics t JOIN users u ON u.id = t.user_id
			WHERE t.status = 'published' AND ($2 IN ('', 'topic'))
			  AND (t.title || ' ' || t.content || ' ' || t.structured_content::text) ILIKE $1
			UNION ALL
			SELECT p.id, 'post'::text, p.topic_id, t.title, LEFT(p.content, 400),
			       u.username, p.created_at, 2
			FROM posts p
			JOIN topics t ON t.id = p.topic_id
			JOIN users u ON u.id = p.user_id
			WHERE p.status = 'published' AND t.status = 'published' AND ($2 IN ('', 'post'))
			  AND p.content ILIKE $1
		)
		SELECT id, type, topic_id, title, excerpt, author_name, created_at
		FROM matched ORDER BY rank, created_at DESC LIMIT $3 OFFSET $4`
	rows, err := d.db.QueryContext(ctx, query, pattern, typeFilter, pageSize, offset)
	if err != nil {
		return items, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		item := &model.SearchResult{}
		if err := rows.Scan(&item.ID, &item.Type, &item.TopicID, &item.Title, &item.Excerpt, &item.AuthorName, &item.CreatedAt); err != nil {
			return items, 0, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return items, 0, err
	}
	countQuery := `
		SELECT
			(SELECT COUNT(*) FROM topics t WHERE t.status='published' AND ($2 IN ('', 'topic')) AND
			 ((t.title || ' ' || t.content || ' ' || t.structured_content::text) ILIKE $1))
			+
			(SELECT COUNT(*) FROM posts p JOIN topics t ON t.id=p.topic_id WHERE p.status='published' AND t.status='published' AND ($2 IN ('', 'post')) AND p.content ILIKE $1)`
	var total int64
	err = d.db.QueryRowContext(ctx, countQuery, pattern, typeFilter).Scan(&total)
	return items, total, err
}
