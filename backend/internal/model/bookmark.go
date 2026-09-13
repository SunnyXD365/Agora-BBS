package model

import "time"

type Bookmark struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	TopicID   int64     `json:"topic_id"`
	CreatedAt time.Time `json:"created_at"`
	Topic     *Topic    `json:"topic,omitempty"`
}
