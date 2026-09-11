package model

import "time"

type Post struct {
	ID            int64      `json:"id" db:"id"`
	TopicID       int64      `json:"topic_id" db:"topic_id"`
	UserID        int64      `json:"user_id" db:"user_id"`
	ParentID      *int64     `json:"parent_id" db:"parent_id"`
	Content       string     `json:"content" db:"content"`
	PostType      string     `json:"post_type" db:"post_type"`
	Status        string     `json:"status" db:"status"`
	CoolingEndsAt *time.Time `json:"cooling_ends_at,omitempty" db:"cooling_ends_at"`
	LikeCount     int        `json:"like_count" db:"like_count"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`

	// 联表查询拓展字段（非物理表字段）
	AuthorName   string `json:"author_name,omitempty" db:"author_name"`
	AuthorAvatar string `json:"author_avatar,omitempty" db:"author_avatar"`
}

// CreatePostReq 发布回复请求 DTO
type CreatePostReq struct {
	TopicID  int64  `json:"topic_id"`
	ParentID *int64 `json:"parent_id"` // 可选：多级回复的父回复 ID
	Content  string `json:"content" binding:"required,min=1"`
	PostType string `json:"post_type" binding:"omitempty,oneof=debate evidence experience thanks"`
}

// PostListReq 回复列表请求 DTO
type PostListReq struct {
	TopicID  int64 `form:"topic_id" binding:"required"`
	Page     int   `form:"page,default=1"`
	PageSize int   `form:"page_size,default=20"`
}
