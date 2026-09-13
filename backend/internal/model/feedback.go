package model

import (
	"encoding/json"
	"time"
)

type Feedback struct {
	ID             int64           `json:"id"`
	UserID         int64           `json:"user_id"`
	TargetType     string          `json:"target_type"`
	TargetID       int64           `json:"target_id"`
	Stance         string          `json:"stance"`
	Tag            string          `json:"tag"`
	Reason         string          `json:"reason"`
	Status         string          `json:"status"`
	LLMAuditStatus string          `json:"llm_audit_status"`
	LLMAuditResult json.RawMessage `json:"llm_audit_result"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type UpsertFeedbackReq struct {
	TargetType string `json:"target_type" binding:"required,oneof=topic post"`
	TargetID   int64  `json:"target_id" binding:"required"`
	Stance     string `json:"stance" binding:"required,oneof=support challenge"`
	Tag        string `json:"tag" binding:"required,oneof=logical new_perspective well_sourced empathetic factual_concern reasoning_gap inappropriate"`
	Reason     string `json:"reason" binding:"required,min=5,max=120"`
}

type FeedbackSummary struct {
	Support   int            `json:"support"`
	Challenge int            `json:"challenge"`
	Score     int            `json:"score"`
	Tags      map[string]int `json:"tags"`
	Mine      *Feedback      `json:"mine,omitempty"`
}

type CommentCluster struct {
	ID         int64     `json:"id"`
	TopicID    int64     `json:"topic_id"`
	Tag        string    `json:"tag"`
	Summary    string    `json:"summary"`
	Weight     float64   `json:"weight"`
	PostIDs    []int64   `json:"post_ids"`
	Generation int       `json:"generation"`
	CreatedAt  time.Time `json:"created_at"`
}
