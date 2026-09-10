package service

import (
	"context"
	"errors"

	"Agora-BBS/internal/dao"
	"Agora-BBS/internal/model"
)

type TopicService struct {
	topicDAO    *dao.TopicDAO
	categoryDAO *dao.CategoryDAO
}

func NewTopicService(topicDAO *dao.TopicDAO, categoryDAO *dao.CategoryDAO) *TopicService {
	return &TopicService{
		topicDAO:    topicDAO,
		categoryDAO: categoryDAO,
	}
}

func (s *TopicService) CreateTopic(ctx context.Context, userID int64, req *model.CreateTopicReq) (*model.Topic, error) {
	// 校验板块是否存在
	cat, err := s.categoryDAO.GetByID(ctx, req.CategoryID)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, errors.New("category not found")
	}

	t := &model.Topic{
		CategoryID: req.CategoryID,
		UserID:     userID,
		Title:      req.Title,
		Content:    req.Content,
		Status:     "normal",
	}
	if err := s.topicDAO.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TopicService) ListTopics(ctx context.Context, req *model.TopicListReq) ([]*model.Topic, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}
	offset := (req.Page - 1) * req.PageSize
	return s.topicDAO.List(ctx, req.CategoryID, offset, req.PageSize)
}

func (s *TopicService) GetTopicDetail(ctx context.Context, id int64) (*model.Topic, error) {
	t, err := s.topicDAO.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, errors.New("topic not found")
	}

	// 异步更新浏览量
	go func() {
		_ = s.topicDAO.IncrementViewCount(context.Background(), id)
	}()

	return t, nil
}
