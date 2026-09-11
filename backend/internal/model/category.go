package model

import "time"

// Category 对应 categories 表
type Category struct {
	ID          int       `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Slug        string    `db:"slug" json:"slug"`
	Description string    `db:"description" json:"description"`
	SortOrder   int       `db:"sort_order" json:"sort_order"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// CreateCategoryReq 创建板块请求 DTO
type CreateCategoryReq struct {
	Name        string `json:"name" binding:"required,min=2,max=32"`
	Slug        string `json:"slug" binding:"required,min=2,max=32"`
	Description string `json:"description" binding:"max=255"`
	SortOrder   int    `json:"sort_order"`
}