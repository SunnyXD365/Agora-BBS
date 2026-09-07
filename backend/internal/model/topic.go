package model

import "time"

type Topic struct {
	ID                int64     `gorm:"primaryKey;column:id" json:"id"`
	CategoryID        int32     `gorm:"column:category_id;not null" json:"category_id"`
	AuthorID          int64     `gorm:"column:author_id;not null" json:"author_id"`
	Title             string    `gorm:"column:title;not null" json:"title"`
	Content           string    `gorm:"column:content;not null" json:"content"`
	StructuredContent string    `gorm:"column:structured_content;type:jsonb;default:'{}'" json:"structured_content"`
	Status            string    `gorm:"column:status;default:'published'" json:"status"`
	ViewCount         int32     `gorm:"column:view_count;default:0" json:"view_count"`
	ReplyCount        int32     `gorm:"column:reply_count;default:0" json:"reply_count"`
	CreatedAt         time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Topic) TableName() string {
	return "topics"
}
