package service

import (
	"context"
	"errors"

	"Agora-BBS/internal/dao"
	"Agora-BBS/internal/model"
)

type LikeService struct {
	likeDAO  *dao.LikeDAO
	topicDAO *dao.TopicDAO
	postDAO  *dao.PostDAO
}

func NewLikeService(likeDAO *dao.LikeDAO, topicDAO *dao.TopicDAO, postDAO *dao.PostDAO) *LikeService {
	return &LikeService{
		likeDAO:  likeDAO,
		topicDAO: topicDAO,
		postDAO:  postDAO,
	}
}

func (s *LikeService) ToggleLike(ctx context.Context, userID int64, req *model.ToggleLikeReq) (*model.ToggleLikeResp, error) {
	// 1. 检查是否已经点过赞
	existingLike, err := s.likeDAO.GetLike(ctx, userID, req.TargetType, req.TargetID)
	if err != nil {
		return nil, err
	}

	isLiked := false
	delta := 0

	if existingLike == nil {
		// 未点赞 -> 添加点赞
		if err := s.likeDAO.CreateLike(ctx, userID, req.TargetType, req.TargetID); err != nil {
			return nil, err
		}
		isLiked = true
		delta = 1
	} else {
		// 已点赞 -> 取消点赞
		if err := s.likeDAO.DeleteLike(ctx, userID, req.TargetType, req.TargetID); err != nil {
			return nil, err
		}
		isLiked = false
		delta = -1
	}

	// 2. 更新对应目标的点赞数
	var newCount int
	if req.TargetType == "topic" {
		newCount, err = s.topicDAO.UpdateLikeCount(ctx, req.TargetID, delta)
	} else if req.TargetType == "post" {
		newCount, err = s.postDAO.UpdateLikeCount(ctx, req.TargetID, delta)
	} else {
		return nil, errors.New("invalid target type")
	}

	if err != nil {
		return nil, err
	}

	return &model.ToggleLikeResp{
		IsLiked:   isLiked,
		LikeCount: newCount,
	}, nil
}
