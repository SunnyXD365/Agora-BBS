package service

import (
	"context"
	"encoding/json"
	"errors"

	"agora-backend/internal/dao"
	"agora-backend/internal/model"
)

type TopicService struct {
	topicDAO *dao.TopicDAO
}

func NewTopicService(topicDAO *dao.TopicDAO) *TopicService {
	return &TopicService{topicDAO: topicDAO}
}

func (s *TopicService) CreateTopic(ctx context.Context, userID int64, req *model.CreateTopicReq) (int64, error) {
	content := req.Content
	structured := req.StructuredContent
	if structured == nil {
		if content == "" {
			return 0, errors.New("content or structured_content is required")
		}
		structured = &model.StructuredContent{Claim: content}
	} else {
		content = structured.Claim
	}
	encoded, err := json.Marshal(structured)
	if err != nil {
		return 0, err
	}
	topic := &model.Topic{
		CategoryID:        req.CategoryID,
		UserID:            userID,
		Title:             req.Title,
		Content:           req.Content,
		StructuredContent: encoded,
	}
	topic.Content = content
	if err := s.topicDAO.CreateTopic(ctx, topic); err != nil {
		return 0, err
	}
	return topic.ID, nil
}

func (s *TopicService) GetTopicDetail(ctx context.Context, topicID int64) (*model.Topic, error) {
	if err := s.topicDAO.IncrementViewCount(ctx, topicID); err != nil {
		return nil, err
	}
	topic, err := s.topicDAO.GetTopicByID(ctx, topicID)
	if err != nil {
		return nil, err
	}

	return topic, nil
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
