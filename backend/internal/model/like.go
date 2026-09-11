package model

import "time"

// Like 对应 likes 表
type Like struct {
	ID         int64     `db:"id" json:"id"`
	UserID     int64     `db:"user_id" json:"user_id"`
	TargetType string    `db:"target_type" json:"target_type"` // 'topic' 或 'post'
	TargetID   int64     `db:"target_id" json:"target_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

// ToggleLikeReq 点赞/取消点赞请求 DTO
type ToggleLikeReq struct {
	TargetType string `json:"target_type" binding:"required,oneof=topic post"`
	TargetID   int64  `json:"target_id" binding:"required"`
}

// ToggleLikeResp 点赞切换响应 DTO
type ToggleLikeResp struct {
	IsLiked   bool `json:"is_liked"`
	LikeCount int  `json:"like_count"`
}