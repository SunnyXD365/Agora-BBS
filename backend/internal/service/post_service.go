package service

import (
	"context"
	"errors"

	"Agora-BBS/internal/dao"
	"Agora-BBS/internal/model"
)

type PostService struct {
	postDAO  *dao.PostDAO
	topicDAO *dao.TopicDAO
}

func NewPostService(postDAO *dao.PostDAO, topicDAO *dao.TopicDAO) *PostService {
	return &PostService{
		postDAO:  postDAO,
		topicDAO: topicDAO,
	}
}

func (s *PostService) CreatePost(ctx context.Context, userID int64, req *model.CreatePostReq) (*model.Post, error) {
	// 校验帖子是否存在
	topic, err := s.topicDAO.GetByID(ctx, req.TopicID)
	if err != nil {
		return nil, err
	}
	if topic == nil {
		return nil, errors.New("topic not found")
	}

	p := &model.Post{
		TopicID:  req.TopicID,
		UserID:   userID,
		ParentID: req.ParentID,
		Content:  req.Content,
		Status:   "normal",
	}

	if err := s.postDAO.Create(ctx, p); err != nil {
		return nil, err
	}

	// 增加主帖的 post_count
	_ = s.topicDAO.IncrementPostCount(ctx, req.TopicID)

	return p, nil
}

func (s *PostService) ListPosts(ctx context.Context, req *model.PostListReq) ([]*model.Post, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}
	offset := (req.Page - 1) * req.PageSize
	return s.postDAO.ListByTopic(ctx, req.TopicID, offset, req.PageSize)
}
