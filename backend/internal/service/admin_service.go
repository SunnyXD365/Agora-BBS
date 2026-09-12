package service

import (
	"context"
	"errors"

	"agora-backend/internal/cache"
	"agora-backend/internal/dao"
	"agora-backend/internal/model"
	"agora-backend/internal/workflow"
)

var ErrSelfSuspend = errors.New("administrator cannot suspend the current account")

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
func (s *AdminService) Overview(ctx context.Context, days int) (*model.AdminOverview, error) {
	if days != 7 && days != 30 && days != 90 {
		days = 7
	}
	return s.dao.Overview(ctx, days)
}
func (s *AdminService) Users(ctx context.Context, page, pageSize int) (*model.Page[*model.AdminUser], error) {
	page, pageSize = normalizePage(page, pageSize)
	items, total, err := s.dao.ListUsers(ctx, page, pageSize)
	return &model.Page[*model.AdminUser]{Items: items, Total: total, Page: page, PageSize: pageSize}, err
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
func (s *AdminService) Contents(ctx context.Context, kind, status string, page, pageSize int) (*model.Page[*model.AdminContent], error) {
	page, pageSize = normalizePage(page, pageSize)
	items, total, err := s.dao.ListContent(ctx, kind, status, page, pageSize)
	return &model.Page[*model.AdminContent]{Items: items, Total: total, Page: page, PageSize: pageSize}, err
}
func (s *AdminService) SetContentVisibility(ctx context.Context, kind string, id int64, hidden bool) error {
	err := s.dao.SetContentVisibility(ctx, kind, id, hidden)
	if err == nil {
		_ = s.cache.DeletePrefix(ctx, "topics:")
	}
	return err
}
func (s *AdminService) LLMJobs(ctx context.Context, status string, page, pageSize int) (*model.Page[*model.AdminLLMJob], error) {
	page, pageSize = normalizePage(page, pageSize)
	items, total, err := s.dao.ListLLMJobs(ctx, status, page, pageSize)
	return &model.Page[*model.AdminLLMJob]{Items: items, Total: total, Page: page, PageSize: pageSize}, err
}
func (s *AdminService) TrustLogs(ctx context.Context, userID int64, page, pageSize int) (*model.Page[*model.AdminTrustLog], error) {
	page, pageSize = normalizePage(page, pageSize)
	items, total, err := s.dao.ListTrustLogs(ctx, userID, page, pageSize)
	return &model.Page[*model.AdminTrustLog]{Items: items, Total: total, Page: page, PageSize: pageSize}, err
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
func (s *AdminService) Categories(ctx context.Context) ([]*model.Category, error) {
	return s.dao.ListAllCategories(ctx)
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
