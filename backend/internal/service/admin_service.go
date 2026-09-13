package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"agora-backend/internal/cache"
	"agora-backend/internal/dao"
	"agora-backend/internal/model"
	"agora-backend/internal/workflow"
)

var ErrSelfSuspend = errors.New("administrator cannot suspend the current account")
var ErrSelfDemote = errors.New("administrator cannot remove the current account's administrator role")

type AdminService struct {
	dao     *dao.AdminDAO
	cache   *cache.Store
	llm     workflow.LLMStarter
	reviews workflow.BlindReviewStarter
}

func NewAdminService(adminDAO *dao.AdminDAO, store *cache.Store, llm workflow.LLMStarter, reviews workflow.BlindReviewStarter) *AdminService {
	return &AdminService{dao: adminDAO, cache: store, llm: llm, reviews: reviews}
}

func normalizePage(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return page, pageSize
}

func normalizeAdminPage(query *model.AdminPageQuery, allowedSort map[string]bool, defaultSort string) {
	query.Page, query.PageSize = normalizePage(query.Page, query.PageSize)
	query.Query = strings.TrimSpace(query.Query)
	if runes := []rune(query.Query); len(runes) > 100 {
		query.Query = string(runes[:100])
	}
	if !allowedSort[query.Sort] {
		query.Sort = defaultSort
	}
	query.Order = strings.ToLower(query.Order)
	if query.Order != "asc" && query.Order != "desc" {
		query.Order = "desc"
	}
}
func (s *AdminService) Overview(ctx context.Context, days int) (*model.AdminOverview, error) {
	if days != 7 && days != 30 && days != 90 {
		days = 7
	}
	return s.dao.Overview(ctx, days)
}
func (s *AdminService) Users(ctx context.Context, query *model.AdminUserQuery) (*model.Page[*model.AdminUser], error) {
	normalizeAdminPage(&query.AdminPageQuery, map[string]bool{"created_at": true, "username": true, "level": true, "trust_score": true, "reading_seconds": true}, "created_at")
	if query.Level < 0 || query.Level > 3 {
		query.Level = -1
	}
	if query.Role != "user" && query.Role != "admin" {
		query.Role = ""
	}
	if query.Status != "active" && query.Status != "suspended" {
		query.Status = ""
	}
	items, total, err := s.dao.ListUsers(ctx, query)
	return &model.Page[*model.AdminUser]{Items: items, Total: total, Page: query.Page, PageSize: query.PageSize}, err
}
func (s *AdminService) SetUserStatus(ctx context.Context, actorID, targetID int64, status string) error {
	if actorID == targetID && status == "suspended" {
		return ErrSelfSuspend
	}
	if status != "active" && status != "suspended" {
		return errors.New("invalid user status")
	}
	return s.dao.SetUserStatus(ctx, targetID, status)
}
func (s *AdminService) UpdateUser(ctx context.Context, actorID, targetID int64, req *model.UpdateAdminUserReq) (*model.AdminUser, error) {
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.Reason = strings.TrimSpace(req.Reason)
	usernameLength := len([]rune(req.Username))
	reasonLength := len([]rune(req.Reason))
	if usernameLength < 3 || usernameLength > 32 {
		return nil, errors.New("username must contain 3 to 32 characters")
	}
	if reasonLength < 5 || reasonLength > 200 {
		return nil, errors.New("adjustment reason must contain 5 to 200 characters")
	}
	if req.Role != "user" && req.Role != "admin" {
		return nil, errors.New("invalid user role")
	}
	if req.Role == "admin" && req.Email == "" {
		return nil, errors.New("administrator accounts must have a verification email")
	}
	if req.Status != "active" && req.Status != "suspended" {
		return nil, errors.New("invalid user status")
	}
	if req.UnlockLevel < 0 || req.UnlockLevel > 3 || req.TrustScore < -1000 || req.TrustScore > 10000 {
		return nil, errors.New("level or trust score is outside the allowed range")
	}
	if actorID == targetID && req.Status == "suspended" {
		return nil, ErrSelfSuspend
	}
	if actorID == targetID && req.Role != "admin" {
		return nil, ErrSelfDemote
	}
	item, err := s.dao.UpdateUser(ctx, actorID, targetID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return item, nil
}
func (s *AdminService) Contents(ctx context.Context, query *model.AdminContentQuery) (*model.Page[*model.AdminContent], error) {
	normalizeAdminPage(&query.AdminPageQuery, map[string]bool{"created_at": true, "title": true, "author": true, "status": true}, "created_at")
	if query.Kind != "topic" && query.Kind != "post" {
		return nil, errors.New("invalid content type")
	}
	items, total, err := s.dao.ListContent(ctx, query)
	return &model.Page[*model.AdminContent]{Items: items, Total: total, Page: query.Page, PageSize: query.PageSize}, err
}
func (s *AdminService) SetContentVisibility(ctx context.Context, kind string, id int64, hidden bool) error {
	err := s.dao.SetContentVisibility(ctx, kind, id, hidden)
	if err == nil {
		_ = s.cache.DeletePrefix(ctx, "topics:")
	}
	return err
}
func (s *AdminService) LLMJobs(ctx context.Context, query *model.AdminLLMJobQuery) (*model.Page[*model.AdminLLMJob], error) {
	normalizeAdminPage(&query.AdminPageQuery, map[string]bool{"created_at": true, "latency_ms": true, "tokens": true, "attempts": true, "status": true}, "created_at")
	items, total, err := s.dao.ListLLMJobs(ctx, query)
	return &model.Page[*model.AdminLLMJob]{Items: items, Total: total, Page: query.Page, PageSize: query.PageSize}, err
}
func (s *AdminService) TrustLogs(ctx context.Context, query *model.AdminTrustLogQuery) (*model.Page[*model.AdminTrustLog], error) {
	normalizeAdminPage(&query.AdminPageQuery, map[string]bool{"created_at": true, "score_delta": true, "username": true, "event_type": true}, "created_at")
	items, total, err := s.dao.ListTrustLogs(ctx, query)
	return &model.Page[*model.AdminTrustLog]{Items: items, Total: total, Page: query.Page, PageSize: query.PageSize}, err
}
func (s *AdminService) RetryLLMJob(ctx context.Context, id int64) error {
	job, err := s.dao.GetLLMJob(ctx, id)
	if err != nil {
		return err
	}
	if job.Status != "failed" {
		return errors.New("only failed jobs can be retried")
	}
	switch job.JobType {
	case "feedback_audit":
		return s.llm.StartFeedbackAudit(ctx, job.AggregateID)
	case "review_audit":
		return s.reviews.StartReviewAudit(ctx, job.AggregateID)
	case "review_fallback":
		return s.reviews.StartFinalizeReview(ctx, job.AggregateID)
	case "comment_cluster":
		return s.llm.StartClusterTopic(ctx, job.AggregateID)
	default:
		return errors.New("unsupported job type")
	}
}
func (s *AdminService) Categories(ctx context.Context, query *model.AdminCategoryQuery) ([]*model.Category, error) {
	query.Query = strings.TrimSpace(query.Query)
	if runes := []rune(query.Query); len(runes) > 100 {
		query.Query = string(runes[:100])
	}
	allowed := map[string]bool{"sort_order": true, "name": true, "created_at": true, "status": true}
	if !allowed[query.Sort] {
		query.Sort = "sort_order"
	}
	query.Order = strings.ToLower(query.Order)
	if query.Order != "asc" && query.Order != "desc" {
		query.Order = "asc"
	}
	return s.dao.ListAllCategories(ctx, query)
}
func (s *AdminService) CreateCategory(ctx context.Context, req *model.UpdateCategoryReq) (*model.Category, error) {
	item, err := s.dao.CreateCategory(ctx, req)
	if err == nil {
		_ = s.cache.DeletePrefix(ctx, "categories:")
	}
	return item, err
}
func (s *AdminService) UpdateCategory(ctx context.Context, id int64, req *model.UpdateCategoryReq) (*model.Category, error) {
	item, err := s.dao.UpdateCategory(ctx, id, req)
	if err == nil {
		_ = s.cache.DeletePrefix(ctx, "categories:")
		_ = s.cache.DeletePrefix(ctx, "topics:")
	}
	return item, err
}
