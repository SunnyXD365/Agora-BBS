package service

import (
	"context"
	"errors"

	"agora-backend/internal/dao"
	"agora-backend/internal/model"
)

type LikeService struct {
	likeDAO *dao.LikeDAO
}

func NewLikeService(likeDAO *dao.LikeDAO) *LikeService {
	return &LikeService{likeDAO: likeDAO}
}

// Like 对主题帖或回复点赞
func (s *LikeService) Like(ctx context.Context, userID int64, req *model.ToggleLikeReq) error {
	if req.TargetType != "topic" && req.TargetType != "post" {
		return errors.New("invalid target_type, must be 'topic' or 'post'")
	}

	like := &model.Like{
		UserID:     userID,
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
	}

	return s.likeDAO.CreateLike(ctx, like)
}

// Unlike 取消点赞
func (s *LikeService) Unlike(ctx context.Context, userID int64, req *model.ToggleLikeReq) error {
	if req.TargetType != "topic" && req.TargetType != "post" {
		return errors.New("invalid target_type, must be 'topic' or 'post'")
	}

	like := &model.Like{
		UserID:     userID,
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
	}

	return s.likeDAO.DeleteLike(ctx, like)
}
