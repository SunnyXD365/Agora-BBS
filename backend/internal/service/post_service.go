package service

import (
	"context"
	"time"

	"agora-backend/internal/dao"
	"agora-backend/internal/model"
	"agora-backend/internal/workflow"
)

type PostService struct {
	postDAO    *dao.PostDAO
	governance *GovernanceService
	cooling    workflow.CoolingStarter
}

func NewPostService(postDAO *dao.PostDAO, governance *GovernanceService, cooling workflow.CoolingStarter) *PostService {
	return &PostService{postDAO: postDAO, governance: governance, cooling: cooling}
}

func (s *PostService) CreatePost(ctx context.Context, userID int64, req *model.CreatePostReq) (*model.Post, error) {
	if err := s.governance.RequireReplyPermission(ctx, userID, req.TopicID); err != nil {
		return nil, err
	}
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
	endsAt := time.Now().Add(time.Duration(s.governance.Policy().CoolingSeconds) * time.Second)
	post.Status = "cooling"
	post.CoolingEndsAt = &endsAt
	if err := s.postDAO.CreatePost(ctx, post); err != nil {
		return nil, err
	}
	if err := s.cooling.StartCooling(ctx, workflow.CoolingInput{Kind: "post", ID: post.ID, EndsAt: endsAt}); err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostService) ListPosts(ctx context.Context, topicID, viewerID int64, page, pageSize int) (*model.Page[*model.Post], error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 50 {
		pageSize = 20
	}

	// 保持参数位置正确：(ctx, topicID, page, pageSize)
	posts, total, err := s.postDAO.ListPostsByTopicID(ctx, topicID, viewerID, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &model.Page[*model.Post]{Items: posts, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *PostService) UpdateCooling(ctx context.Context, userID, postID int64, req *model.UpdatePostReq) (*model.Post, error) {
	endsAt := time.Now().Add(time.Duration(s.governance.Policy().CoolingSeconds) * time.Second)
	post := &model.Post{ID: postID, UserID: userID, Content: req.Content, PostType: req.PostType, Status: "cooling", CoolingEndsAt: &endsAt}
	if err := s.postDAO.UpdateCooling(ctx, post); err != nil {
		return nil, err
	}
	if err := s.cooling.StartCooling(ctx, workflow.CoolingInput{Kind: "post", ID: postID, EndsAt: endsAt}); err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostService) RecallCooling(ctx context.Context, userID, postID int64) error {
	return s.postDAO.RecallCooling(ctx, postID, userID)
}
