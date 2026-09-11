package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	"unicode/utf8"

	"agora-backend/internal/dao"
	"agora-backend/internal/model"
	"agora-backend/internal/workflow"
)

type TopicService struct {
	topicDAO   *dao.TopicDAO
	governance *GovernanceService
	cooling    workflow.CoolingStarter
}

func NewTopicService(topicDAO *dao.TopicDAO, governance *GovernanceService, cooling workflow.CoolingStarter) *TopicService {
	return &TopicService{topicDAO: topicDAO, governance: governance, cooling: cooling}
}

func (s *TopicService) CreateTopic(ctx context.Context, userID int64, req *model.CreateTopicReq) (*model.Topic, error) {
	if err := s.governance.RequireTopicPermission(ctx, userID); err != nil {
		return nil, err
	}
	content := req.Content
	structured := req.StructuredContent
	if structured == nil {
		if content == "" {
			return nil, errors.New("content or structured_content is required")
		}
		structured = &model.StructuredContent{Claim: content}
	} else {
		content = structured.Claim
	}
	encoded, err := json.Marshal(structured)
	if err != nil {
		return nil, err
	}
	categoryRequiresReview, err := s.topicDAO.CategoryRequiresReview(ctx, req.CategoryID)
	if err != nil {
		return nil, err
	}
	policy := s.governance.Policy()
	endsAt := time.Now().Add(time.Duration(policy.CoolingSeconds) * time.Second)
	topic := &model.Topic{
		CategoryID:        req.CategoryID,
		UserID:            userID,
		Title:             req.Title,
		Content:           req.Content,
		StructuredContent: encoded,
		Status:            "cooling",
		CoolingEndsAt:     &endsAt,
	}
	topic.Content = content
	if err := s.topicDAO.CreateTopic(ctx, topic); err != nil {
		return nil, err
	}
	requiresReview := categoryRequiresReview || utf8.RuneCountInString(structured.Claim+structured.Evidence+structured.Uncertainty) >= policy.LongTopicChars
	if err := s.cooling.StartCooling(ctx, workflow.CoolingInput{Kind: "topic", ID: topic.ID, EndsAt: endsAt, RequiresReview: requiresReview}); err != nil {
		return nil, err
	}
	return topic, nil
}

func (s *TopicService) GetTopicDetail(ctx context.Context, topicID, viewerID int64) (*model.Topic, error) {
	if err := s.topicDAO.IncrementViewCount(ctx, topicID); err != nil {
		return nil, err
	}
	topic, err := s.topicDAO.GetTopicByID(ctx, topicID, viewerID)
	if err != nil {
		return nil, err
	}

	return topic, nil
}

func (s *TopicService) UpdateCooling(ctx context.Context, userID, topicID int64, req *model.UpdateTopicReq) (*model.Topic, error) {
	encoded, err := json.Marshal(req.StructuredContent)
	if err != nil {
		return nil, err
	}
	policy := s.governance.Policy()
	endsAt := time.Now().Add(time.Duration(policy.CoolingSeconds) * time.Second)
	topic := &model.Topic{ID: topicID, UserID: userID, Title: req.Title, Content: req.StructuredContent.Claim, StructuredContent: encoded, Status: "cooling", CoolingEndsAt: &endsAt}
	if err := s.topicDAO.UpdateCooling(ctx, topic); err != nil {
		return nil, err
	}
	requiresReview := utf8.RuneCountInString(req.StructuredContent.Claim+req.StructuredContent.Evidence+req.StructuredContent.Uncertainty) >= policy.LongTopicChars
	if err := s.cooling.StartCooling(ctx, workflow.CoolingInput{Kind: "topic", ID: topicID, EndsAt: endsAt, RequiresReview: requiresReview}); err != nil {
		return nil, err
	}
	return topic, nil
}

func (s *TopicService) RecallCooling(ctx context.Context, userID, topicID int64) error {
	return s.topicDAO.RecallCooling(ctx, topicID, userID)
}

func (s *TopicService) ListTopics(ctx context.Context, req *model.TopicListReq) (*model.Page[*model.Topic], error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 50 {
		req.PageSize = 20
	}

	// 保持参数位置正确：(ctx, categoryID, page, pageSize)
	topics, total, err := s.topicDAO.ListTopicsByCategoryID(ctx, req.CategoryID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	return &model.Page[*model.Topic]{Items: topics, Total: total, Page: req.Page, PageSize: req.PageSize}, nil
}
