package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"unicode/utf8"

	"agora-backend/internal/cache"
	"agora-backend/internal/dao"
	"agora-backend/internal/model"
	"agora-backend/internal/workflow"
)

type TopicService struct {
	topicDAO   *dao.TopicDAO
	governance *GovernanceService
	cooling    workflow.CoolingStarter
	cache      *cache.Store
}

func NewTopicService(topicDAO *dao.TopicDAO, governance *GovernanceService, cooling workflow.CoolingStarter, store *cache.Store) *TopicService {
	return &TopicService{topicDAO: topicDAO, governance: governance, cooling: cooling, cache: store}
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
	_ = s.cache.DeletePrefix(ctx, "topics:list:")
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
	key := fmt.Sprintf("topics:detail:%d:%d", topicID, viewerID)
	var topic model.Topic
	if err := s.cache.GetJSON(ctx, key, &topic); err == nil {
		return &topic, nil
	}
	loaded, err := s.topicDAO.GetTopicByID(ctx, topicID, viewerID)
	if err != nil {
		return nil, err
	}
	if loaded != nil {
		_ = s.cache.SetJSON(ctx, key, loaded, 30*time.Second)
	}
	return loaded, nil
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
	_ = s.cache.DeletePrefix(ctx, fmt.Sprintf("topics:detail:%d:", topicID))
	requiresReview := utf8.RuneCountInString(req.StructuredContent.Claim+req.StructuredContent.Evidence+req.StructuredContent.Uncertainty) >= policy.LongTopicChars
	if err := s.cooling.StartCooling(ctx, workflow.CoolingInput{Kind: "topic", ID: topicID, EndsAt: endsAt, RequiresReview: requiresReview}); err != nil {
		return nil, err
	}
	return topic, nil
}

func (s *TopicService) RecallCooling(ctx context.Context, userID, topicID int64) error {
	err := s.topicDAO.RecallCooling(ctx, topicID, userID)
	if err == nil {
		_ = s.cache.DeletePrefix(ctx, fmt.Sprintf("topics:detail:%d:", topicID))
		_ = s.cache.DeletePrefix(ctx, "topics:list:")
	}
	return err
}

func (s *TopicService) ListTopics(ctx context.Context, req *model.TopicListReq) (*model.Page[*model.Topic], error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 50 {
		req.PageSize = 20
	}

	// 保持参数位置正确：(ctx, categoryID, page, pageSize)
	key := fmt.Sprintf("topics:list:%d:%d:%d", req.CategoryID, req.Page, req.PageSize)
	var cached model.Page[*model.Topic]
	if err := s.cache.GetJSON(ctx, key, &cached); err == nil {
		return &cached, nil
	}
	topics, total, err := s.topicDAO.ListTopicsByCategoryID(ctx, req.CategoryID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	result := &model.Page[*model.Topic]{Items: topics, Total: total, Page: req.Page, PageSize: req.PageSize}
	_ = s.cache.SetJSON(ctx, key, result, 30*time.Second)
	return result, nil
}
