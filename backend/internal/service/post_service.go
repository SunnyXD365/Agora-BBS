package service

import (
	"context"

	"agora-backend/internal/dao"
	"agora-backend/internal/model"
)

type PostService struct {
	postDAO *dao.PostDAO
}

func NewPostService(postDAO *dao.PostDAO) *PostService {
	return &PostService{postDAO: postDAO}
}

func (s *PostService) CreatePost(ctx context.Context, userID int64, req *model.CreatePostReq) (int64, error) {
	if req.PostType == "" {
		req.PostType = "experience"
	}
	post := &model.Post{
		TopicID:  req.TopicID,
		UserID:   userID,
		ParentID: req.ParentID,
		Content:  req.Content,
		PostType: req.PostType,
	}
	if err := s.postDAO.CreatePost(ctx, post); err != nil {
		return 0, err
	}
	return post.ID, nil
}

func (s *PostService) ListPosts(ctx context.Context, topicID int64, page, pageSize int) (*model.Page[*model.Post], error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 50 {
		pageSize = 20
	}

	// 保持参数位置正确：(ctx, topicID, page, pageSize)
	posts, total, err := s.postDAO.ListPostsByTopicID(ctx, topicID, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &model.Page[*model.Post]{Items: posts, Total: total, Page: page, PageSize: pageSize}, nil
}
