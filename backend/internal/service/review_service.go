package service

import (
	"context"

	"agora-backend/internal/dao"
	"agora-backend/internal/model"
	"agora-backend/internal/workflow"
)

type ReviewService struct {
	dao        *dao.ReviewDAO
	governance *GovernanceService
	starter    workflow.BlindReviewStarter
}

func NewReviewService(reviewDAO *dao.ReviewDAO, governance *GovernanceService, starter workflow.BlindReviewStarter) *ReviewService {
	return &ReviewService{dao: reviewDAO, governance: governance, starter: starter}
}
func (s *ReviewService) List(ctx context.Context, userID int64) ([]*model.ReviewTask, error) {
	if err := s.governance.RequireReviewPermission(ctx, userID); err != nil {
		return nil, err
	}
	return s.dao.ListTasks(ctx, userID)
}
func (s *ReviewService) Submit(ctx context.Context, userID, taskID int64, req *model.SubmitReviewReq) (*model.ReviewSubmission, error) {
	if err := s.governance.RequireReviewPermission(ctx, userID); err != nil {
		return nil, err
	}
	result, err := s.dao.Submit(ctx, userID, taskID, req)
	if err == nil {
		_ = s.starter.StartReviewAudit(ctx, taskID)
	}
	return result, err
}
