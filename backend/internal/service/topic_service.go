package service

import (
	"context"

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
	topic := &model.Topic{
		CategoryID: req.CategoryID,
		UserID:     userID,
		Title:      req.Title,
		Content:    req.Content,
	}
	if err := s.topicDAO.CreateTopic(ctx, topic); err != nil {
		return 0, err
	}
	return topic.ID, nil
}

func (s *TopicService) GetTopicDetail(ctx context.Context, topicID int64) (*model.Topic, error) {
	topic, err := s.topicDAO.GetTopicByID(ctx, topicID)
	if err != nil {
		return nil, err
	}

	// 触发浏览量自增
	_ = s.topicDAO.IncrementViewCount(ctx, topicID)
	return topic, nil
}

func (s *TopicService) ListTopics(ctx context.Context, req *model.TopicListReq) ([]*model.Topic, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 50 {
		req.PageSize = 20
	}

	// 保持参数位置正确：(ctx, categoryID, page, pageSize)
	topics, err := s.topicDAO.ListTopicsByCategoryID(ctx, req.CategoryID, req.Page, req.PageSize)
	if err != nil {
		return make([]*model.Topic, 0), err
	}
	return topics, nil
}