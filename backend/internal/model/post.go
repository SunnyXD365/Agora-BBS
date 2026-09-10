package model

import "time"

// Post 对应 posts 表 (帖子回复/评论)
type Post struct {
	ID         int64     `db:"id" json:"id"`
	TopicID    int64     `db:"topic_id" json:"topic_id"`
	UserID     int64     `db:"user_id" json:"user_id"`
	ParentID   *int64    `db:"parent_id" json:"parent_id"`
	Content    string    `db:"content" json:"content"`
	LikeCount  int       `db:"like_count" json:"like_count"`
	Status     string    `db:"status" json:"status"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`

	AuthorName   string `db:"author_name" json:"author_name,omitempty"`
	AuthorAvatar string `db:"author_avatar" json:"author_avatar,omitempty"`
}

// CreatePostReq 发布回复 DTO
type CreatePostReq struct {
	TopicID  int64  `json:"topic_id" binding:"required"`
	ParentID *int64 `json:"parent_id"` // 可选：楼中楼回复
	Content  string `json:"content" binding:"required,min=1"`
}

// PostListReq 回复列表筛选 DTO
type PostListReq struct {
	TopicID  int64 `form:"topic_id" binding:"required"`
	Page     int   `form:"page,default=1"`
	PageSize int   `form:"page_size,default=20"`
}
