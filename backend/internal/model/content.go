package model

import "time"

type MyContentItem struct {
	ID            int64      `json:"id"`
	Type          string     `json:"type"`
	TopicID       int64      `json:"topic_id"`
	Title         string     `json:"title"`
	Excerpt       string     `json:"excerpt"`
	Status        string     `json:"status"`
	PostType      string     `json:"post_type,omitempty"`
	CoolingEndsAt *time.Time `json:"cooling_ends_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type SearchResult struct {
	ID         int64     `json:"id"`
	Type       string    `json:"type"`
	TopicID    int64     `json:"topic_id"`
	Title      string    `json:"title"`
	Excerpt    string    `json:"excerpt"`
	AuthorName string    `json:"author_name"`
	CreatedAt  time.Time `json:"created_at"`
}
