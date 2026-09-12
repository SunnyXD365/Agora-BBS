package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"agora-backend/internal/model"
)

type AdminDAO struct{ db *sql.DB }

func NewAdminDAO(db *sql.DB) *AdminDAO { return &AdminDAO{db: db} }

func (d *AdminDAO) Overview(ctx context.Context, days int) (*model.AdminOverview, error) {
	result := &model.AdminOverview{ContentStatus: map[string]int64{}, TrustDistribution: map[string]int64{}, LevelDistribution: map[string]int64{}, FeedbackDistribution: map[string]int64{}, Trend: make([]model.AdminDailyTrend, 0), Categories: make([]model.AdminCategoryStat, 0), TrendDays: days}
	if err := d.db.QueryRowContext(ctx, `SELECT (SELECT COUNT(*) FROM users),(SELECT COUNT(*) FROM topics),(SELECT COUNT(*) FROM posts)`).Scan(&result.UsersTotal, &result.TopicsTotal, &result.PostsTotal); err != nil {
		return nil, err
	}
	if err := d.db.QueryRowContext(ctx, `SELECT COUNT(DISTINCT user_id) FROM (SELECT user_id FROM topics WHERE created_at>=CURRENT_TIMESTAMP-INTERVAL '7 days' UNION ALL SELECT user_id FROM posts WHERE created_at>=CURRENT_TIMESTAMP-INTERVAL '7 days' UNION ALL SELECT user_id FROM contextual_feedbacks WHERE created_at>=CURRENT_TIMESTAMP-INTERVAL '7 days') a`).Scan(&result.ActiveUsers7Days); err != nil {
		return nil, err
	}
	if err := d.db.QueryRowContext(ctx, `SELECT
		(SELECT COUNT(*) FROM users WHERE created_at>=CURRENT_DATE),
		(SELECT COUNT(*) FROM contextual_feedbacks WHERE status='active'),
		(SELECT COUNT(*) FROM bookmarks),
		(SELECT COUNT(*) FROM users WHERE status='suspended'),
		(SELECT COALESCE(SUM(verified_read_seconds),0)/3600.0 FROM user_trust_profiles)`).Scan(&result.NewUsersToday, &result.FeedbackTotal, &result.BookmarksTotal, &result.SuspendedUsers, &result.VerifiedReadHours); err != nil {
		return nil, err
	}
	rows, err := d.db.QueryContext(ctx, `SELECT status,COUNT(*) FROM (SELECT status FROM topics UNION ALL SELECT status FROM posts) c GROUP BY status`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var key string
		var count int64
		if err = rows.Scan(&key, &count); err != nil {
			rows.Close()
			return nil, err
		}
		result.ContentStatus[key] = count
	}
	rows.Close()
	rows, err = d.db.QueryContext(ctx, `SELECT 'L'||unlock_level::text,COUNT(*) FROM user_trust_profiles GROUP BY unlock_level ORDER BY unlock_level`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var key string
		var count int64
		if err = rows.Scan(&key, &count); err != nil {
			rows.Close()
			return nil, err
		}
		result.LevelDistribution[key] = count
	}
	rows.Close()
	rows, err = d.db.QueryContext(ctx, `SELECT tag,COUNT(*) FROM contextual_feedbacks WHERE status='active' GROUP BY tag ORDER BY COUNT(*) DESC`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var key string
		var count int64
		if err = rows.Scan(&key, &count); err != nil {
			rows.Close()
			return nil, err
		}
		result.FeedbackDistribution[key] = count
	}
	rows.Close()
	rows, err = d.db.QueryContext(ctx, `SELECT CASE WHEN trust_score<0 THEN 'negative' WHEN trust_score<20 THEN '0-19' WHEN trust_score<50 THEN '20-49' ELSE '50+' END,COUNT(*) FROM user_trust_profiles GROUP BY 1 ORDER BY 1`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var key string
		var count int64
		if err = rows.Scan(&key, &count); err != nil {
			rows.Close()
			return nil, err
		}
		result.TrustDistribution[key] = count
	}
	rows.Close()
	if err = d.db.QueryRowContext(ctx, `SELECT COUNT(*),COUNT(*) FILTER(WHERE status='completed'),COUNT(*) FILTER(WHERE status='expired') FROM blind_review_batches`).Scan(&result.ReviewTotal, &result.ReviewCompleted, &result.ReviewExpired); err != nil {
		return nil, err
	}
	if err = d.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM blind_review_tasks WHERE llm_check_status='valid'`).Scan(&result.ReviewFair); err != nil {
		return nil, err
	}
	if err = d.db.QueryRowContext(ctx, `SELECT COUNT(*),COUNT(*) FILTER(WHERE status='completed'),COALESCE(AVG(latency_ms) FILTER(WHERE status='completed'),0),COALESCE(SUM(prompt_tokens),0),COALESCE(SUM(completion_tokens),0) FROM llm_jobs`).Scan(&result.LLMCalls, &result.LLMSuccess, &result.LLMAverageMS, &result.LLMPromptTokens, &result.LLMOutputTokens); err != nil {
		return nil, err
	}
	rows, err = d.db.QueryContext(ctx, `SELECT to_char(day,'YYYY-MM-DD'),(SELECT COUNT(*) FROM users WHERE created_at>=day AND created_at<day+INTERVAL '1 day'),(SELECT COUNT(*) FROM topics WHERE created_at>=day AND created_at<day+INTERVAL '1 day'),(SELECT COUNT(*) FROM posts WHERE created_at>=day AND created_at<day+INTERVAL '1 day'),(SELECT COUNT(*) FROM contextual_feedbacks WHERE created_at>=day AND created_at<day+INTERVAL '1 day') FROM generate_series(CURRENT_DATE-($1::int-1),CURRENT_DATE,INTERVAL '1 day') day`, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item model.AdminDailyTrend
		if err = rows.Scan(&item.Date, &item.Users, &item.Topics, &item.Posts, &item.Feedback); err != nil {
			return nil, err
		}
		result.Trend = append(result.Trend, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	rows.Close()
	rows, err = d.db.QueryContext(ctx, `SELECT c.name,COUNT(DISTINCT t.id),COUNT(p.id) FROM categories c LEFT JOIN topics t ON t.category_id=c.id LEFT JOIN posts p ON p.topic_id=t.id GROUP BY c.id,c.name,c.sort_order ORDER BY c.sort_order,c.id`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var item model.AdminCategoryStat
		if err = rows.Scan(&item.Name, &item.Topics, &item.Posts); err != nil {
			return nil, err
		}
		result.Categories = append(result.Categories, item)
	}
	return result, rows.Err()
}

func (d *AdminDAO) ListUsers(ctx context.Context, page, pageSize int) ([]*model.AdminUser, int64, error) {
	rows, err := d.db.QueryContext(ctx, `SELECT u.id,u.username,u.email,u.role,u.status,tp.unlock_level,tp.trust_score,tp.verified_read_seconds,p.onboarding_status,u.created_at FROM users u JOIN user_trust_profiles tp ON tp.user_id=u.id JOIN user_profiles p ON p.user_id=u.id ORDER BY u.created_at DESC LIMIT $1 OFFSET $2`, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]*model.AdminUser, 0)
	for rows.Next() {
		item := &model.AdminUser{}
		if err = rows.Scan(&item.ID, &item.Username, &item.Email, &item.Role, &item.Status, &item.UnlockLevel, &item.TrustScore, &item.VerifiedReadSeconds, &item.OnboardingStatus, &item.CreatedAt); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	var total int64
	if err = d.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&total); err != nil {
		return nil, 0, err
	}
	return items, total, rows.Err()
}

func (d *AdminDAO) SetUserStatus(ctx context.Context, id int64, status string) error {
	result, err := d.db.ExecContext(ctx, `UPDATE users SET status=$2,updated_at=CURRENT_TIMESTAMP WHERE id=$1`, id, status)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (d *AdminDAO) ListContent(ctx context.Context, kind, status string, page, pageSize int) ([]*model.AdminContent, int64, error) {
	if kind != "topic" && kind != "post" {
		return nil, 0, errors.New("invalid content type")
	}
	table, titleExpr := "topics", "title"
	if kind == "post" {
		table, titleExpr = "posts", "'回复 #' || c.id::text"
	}
	filter, args := " WHERE TRUE", []any{pageSize, (page - 1) * pageSize}
	if kind == "topic" {
		filter += " AND c.status <> 'draft'"
	}
	if status != "" {
		filter += " AND c.status=$3"
		args = append(args, status)
	}
	query := fmt.Sprintf(`SELECT c.id,%s,left(c.content,160),u.username,c.status,c.created_at FROM %s c JOIN users u ON u.id=c.user_id%s ORDER BY c.created_at DESC LIMIT $1 OFFSET $2`, titleExpr, table, filter)
	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]*model.AdminContent, 0)
	for rows.Next() {
		item := &model.AdminContent{Type: kind}
		if err = rows.Scan(&item.ID, &item.Title, &item.Excerpt, &item.AuthorName, &item.Status, &item.CreatedAt); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	countFilter := " WHERE TRUE"
	if kind == "topic" {
		countFilter += " AND c.status <> 'draft'"
	}
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s c%s`, table, countFilter)
	countArgs := []any{}
	if status != "" {
		countFilter += " AND c.status=$1"
		countQuery = fmt.Sprintf(`SELECT COUNT(*) FROM %s c%s`, table, countFilter)
		countArgs = []any{status}
	}
	var total int64
	if err = d.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}
	return items, total, rows.Err()
}

func (d *AdminDAO) SetContentVisibility(ctx context.Context, kind string, id int64, hidden bool) error {
	if kind != "topic" && kind != "post" {
		return errors.New("invalid content type")
	}
	status, hiddenAt := "published", "NULL"
	if hidden {
		status, hiddenAt = "hidden", "CURRENT_TIMESTAMP"
	}
	query := fmt.Sprintf(`UPDATE %s SET status=$2,hidden_at=%s,updated_at=CURRENT_TIMESTAMP WHERE id=$1 AND status IN ('published','hidden')`, kind+"s", hiddenAt)
	result, err := d.db.ExecContext(ctx, query, id, status)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (d *AdminDAO) ListLLMJobs(ctx context.Context, status string, page, pageSize int) ([]*model.AdminLLMJob, int64, error) {
	filter, args := "", []any{pageSize, (page - 1) * pageSize}
	if status != "" {
		filter = " WHERE status=$3"
		args = append(args, status)
	}
	rows, err := d.db.QueryContext(ctx, `SELECT id,job_type,aggregate_type,aggregate_id,status,attempts,model,result,error_message,prompt_tokens,completion_tokens,latency_ms,created_at,completed_at FROM llm_jobs`+filter+` ORDER BY created_at DESC LIMIT $1 OFFSET $2`, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]*model.AdminLLMJob, 0)
	for rows.Next() {
		item := &model.AdminLLMJob{}
		if err = rows.Scan(&item.ID, &item.JobType, &item.AggregateType, &item.AggregateID, &item.Status, &item.Attempts, &item.Model, &item.Result, &item.ErrorMessage, &item.PromptTokens, &item.CompletionTokens, &item.LatencyMS, &item.CreatedAt, &item.CompletedAt); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	countQuery := "SELECT COUNT(*) FROM llm_jobs"
	countArgs := []any{}
	if status != "" {
		countQuery += " WHERE status=$1"
		countArgs = []any{status}
	}
	var total int64
	if err = d.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}
	return items, total, rows.Err()
}

func (d *AdminDAO) GetLLMJob(ctx context.Context, id int64) (*model.AdminLLMJob, error) {
	item := &model.AdminLLMJob{}
	err := d.db.QueryRowContext(ctx, `SELECT id,job_type,aggregate_type,aggregate_id,status,attempts,model,result,error_message,prompt_tokens,completion_tokens,latency_ms,created_at,completed_at FROM llm_jobs WHERE id=$1`, id).Scan(&item.ID, &item.JobType, &item.AggregateType, &item.AggregateID, &item.Status, &item.Attempts, &item.Model, &item.Result, &item.ErrorMessage, &item.PromptTokens, &item.CompletionTokens, &item.LatencyMS, &item.CreatedAt, &item.CompletedAt)
	return item, err
}

func (d *AdminDAO) ListTrustLogs(ctx context.Context, userID int64, page, pageSize int) ([]*model.AdminTrustLog, int64, error) {
	filter, args := "", []any{pageSize, (page - 1) * pageSize}
	if userID > 0 {
		filter = " WHERE l.user_id=$3"
		args = append(args, userID)
	}
	rows, err := d.db.QueryContext(ctx, `SELECT l.id,l.user_id,u.username,l.event_type,l.score_delta,l.reason,COALESCE(l.reference_type,''),l.reference_id,l.created_at FROM trust_logs l JOIN users u ON u.id=l.user_id`+filter+` ORDER BY l.created_at DESC LIMIT $1 OFFSET $2`, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]*model.AdminTrustLog, 0)
	for rows.Next() {
		item := &model.AdminTrustLog{}
		if err = rows.Scan(&item.ID, &item.UserID, &item.Username, &item.EventType, &item.ScoreDelta, &item.Reason, &item.ReferenceType, &item.ReferenceID, &item.CreatedAt); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	countQuery, countArgs := "SELECT COUNT(*) FROM trust_logs", []any{}
	if userID > 0 {
		countQuery += " WHERE user_id=$1"
		countArgs = []any{userID}
	}
	var total int64
	if err = d.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}
	return items, total, rows.Err()
}

func (d *AdminDAO) ListAllCategories(ctx context.Context) ([]*model.Category, error) {
	rows, err := d.db.QueryContext(ctx, `SELECT id,name,slug,description,sort_order,is_active,requires_review,created_at,updated_at FROM categories ORDER BY sort_order,id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]*model.Category, 0)
	for rows.Next() {
		item := &model.Category{}
		if err = rows.Scan(&item.ID, &item.Name, &item.Slug, &item.Description, &item.SortOrder, &item.IsActive, &item.RequiresReview, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
func (d *AdminDAO) CreateCategory(ctx context.Context, req *model.UpdateCategoryReq) (*model.Category, error) {
	item := &model.Category{}
	err := d.db.QueryRowContext(ctx, `INSERT INTO categories(name,slug,description,sort_order,is_active,requires_review) VALUES($1,$2,$3,$4,$5,$6) RETURNING id,name,slug,description,sort_order,is_active,requires_review,created_at,updated_at`, strings.TrimSpace(req.Name), strings.TrimSpace(req.Slug), req.Description, req.SortOrder, req.IsActive, req.RequiresReview).Scan(&item.ID, &item.Name, &item.Slug, &item.Description, &item.SortOrder, &item.IsActive, &item.RequiresReview, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}
func (d *AdminDAO) UpdateCategory(ctx context.Context, id int64, req *model.UpdateCategoryReq) (*model.Category, error) {
	item := &model.Category{}
	err := d.db.QueryRowContext(ctx, `UPDATE categories SET name=$2,slug=$3,description=$4,sort_order=$5,is_active=$6,requires_review=$7,updated_at=CURRENT_TIMESTAMP WHERE id=$1 RETURNING id,name,slug,description,sort_order,is_active,requires_review,created_at,updated_at`, id, strings.TrimSpace(req.Name), strings.TrimSpace(req.Slug), req.Description, req.SortOrder, req.IsActive, req.RequiresReview).Scan(&item.ID, &item.Name, &item.Slug, &item.Description, &item.SortOrder, &item.IsActive, &item.RequiresReview, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}
