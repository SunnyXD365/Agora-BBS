package service

import (
	"context"
	"crypto/rand"
	"math/big"

	"agora-backend/internal/dao"
	"agora-backend/internal/model"
	"agora-backend/internal/workflow"
)

type FeedbackService struct {
	dao        *dao.FeedbackDAO
	governance *GovernanceService
	llm        workflow.LLMStarter
}

func NewFeedbackService(feedbackDAO *dao.FeedbackDAO, governance *GovernanceService, llmStarter workflow.LLMStarter) *FeedbackService {
	return &FeedbackService{dao: feedbackDAO, governance: governance, llm: llmStarter}
}
func (s *FeedbackService) Upsert(ctx context.Context, userID int64, req *model.UpsertFeedbackReq) (*model.Feedback, error) {
	if err := s.governance.RequireFeedbackPermission(ctx, userID); err != nil {
		return nil, err
	}
	feedback, err := s.dao.Upsert(ctx, userID, req)
	if err != nil {
		return nil, err
	}
	probability, err := s.dao.AuditProbability(ctx, userID)
	if err != nil {
		return nil, err
	}
	draw, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		return nil, err
	}
	if float64(draw.Int64())/10000 < probability {
		if err := s.dao.QueueAudit(ctx, feedback.ID); err == nil {
			feedback.LLMAuditStatus = "pending"
			_ = s.llm.StartFeedbackAudit(ctx, feedback.ID)
		}
	}
	return feedback, nil
}
func (s *FeedbackService) Withdraw(ctx context.Context, userID, id int64) error {
	return s.dao.Withdraw(ctx, userID, id)
}
func (s *FeedbackService) Summary(ctx context.Context, targetType string, targetID, viewerID int64) (*model.FeedbackSummary, error) {
	return s.dao.Summary(ctx, targetType, targetID, viewerID)
}
func (s *FeedbackService) Clusters(ctx context.Context, topicID int64) ([]*model.CommentCluster, error) {
	return s.dao.ListClusters(ctx, topicID)
}
