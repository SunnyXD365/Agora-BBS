package model

import "time"

// Topic 对应 topics 表
type Topic struct {
	ID         int64     `db:"id" json:"id"`
	CategoryID int64     `db:"category_id" json:"category_id"`
	UserID     int64     `db:"user_id" json:"user_id"`
	Title      string    `db:"title" json:"title"`
	Content    string    `db:"content" json:"content"`
	ViewCount  int       `db:"view_count" json:"view_count"`
	PostCount  int       `db:"post_count" json:"post_count"`
	IsSticky   bool      `db:"is_sticky" json:"is_sticky"`
	IsEssence  bool      `db:"is_essence" json:"is_essence"`
	Status     string    `db:"status" json:"status"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`

	// 关联字段，方便直接展示发帖人信息
	AuthorName   string `db:"author_name" json:"author_name,omitempty"`
	AuthorAvatar string `db:"author_avatar" json:"author_avatar,omitempty"`
}

// CreateTopicReq 发帖 DTO
type CreateTopicReq struct {
	CategoryID int64  `json:"category_id" binding:"required"`
	Title      string `json:"title" binding:"required,min=3,max=128"`
	Content    string `json:"content" binding:"required,min=5"`
}

// TopicListReq 帖子列表筛选/分页 DTO
type TopicListReq struct {
	CategoryID int64 `form:"category_id"`
	Page       int   `form:"page,default=1"`
	PageSize   int   `form:"page_size,default=20"`
}
