package model

import (
	"encoding/json"
	"time"
)

type StructuredContent struct {
	Claim       string `json:"claim" binding:"required,min=5"`
	Evidence    string `json:"evidence"`
	Uncertainty string `json:"uncertainty"`
}

type Topic struct {
	ID                int64           `json:"id" db:"id"`
	CategoryID        int64           `json:"category_id" db:"category_id"`
	UserID            int64           `json:"user_id" db:"user_id"`
	Title             string          `json:"title" db:"title"`
	Content           string          `json:"content" db:"content"`
	StructuredContent json.RawMessage `json:"structured_content" db:"structured_content"`
	Status            string          `json:"status" db:"status"`
	CoolingEndsAt     *time.Time      `json:"cooling_ends_at,omitempty" db:"cooling_ends_at"`
	ViewCount         int             `json:"view_count" db:"view_count"`
	PostCount         int             `json:"post_count" db:"post_count"`
	LikeCount         int             `json:"like_count" db:"like_count"`
	CreatedAt         time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at" db:"updated_at"`

	// 联表查询拓展字段（非物理表字段）
	AuthorName   string `json:"author_name,omitempty" db:"author_name"`
	AuthorAvatar string `json:"author_avatar,omitempty" db:"author_avatar"`
}

// CreateTopicReq 发布主题帖请求 DTO
type CreateTopicReq struct {
	CategoryID        int64              `json:"category_id" binding:"required"`
	Title             string             `json:"title" binding:"required,min=3,max=128"`
	Content           string             `json:"content" binding:"omitempty,min=5"`
	StructuredContent *StructuredContent `json:"structured_content"`
}

// TopicListReq 帖子列表分页筛选 DTO
type TopicListReq struct {
	CategoryID int64 `form:"category_id"`
	Page       int   `form:"page,default=1"`
	PageSize   int   `form:"page_size,default=20"`
}

type UpdateTopicReq struct {
	Title             string            `json:"title" binding:"required,min=3,max=128"`
	StructuredContent StructuredContent `json:"structured_content" binding:"required"`
}

type SaveTopicDraftReq struct {
	CategoryID        int64             `json:"category_id" binding:"required"`
	Title             string            `json:"title" binding:"max=128"`
	StructuredContent DraftTopicContent `json:"structured_content"`
}

type DraftTopicContent struct {
	Claim       string `json:"claim"`
	Evidence    string `json:"evidence"`
	Uncertainty string `json:"uncertainty"`
}
