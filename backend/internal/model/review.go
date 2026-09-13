package model

import (
	"encoding/json"
	"time"
)

type ReviewTask struct {
	ID          int64           `json:"id"`
	BatchID     int64           `json:"batch_id"`
	SubjectType string          `json:"subject_type"`
	Subject     json.RawMessage `json:"subject"`
	TaskStatus  string          `json:"task_status"`
	Deadline    time.Time       `json:"deadline"`
	CreatedAt   time.Time       `json:"created_at"`
}

type SubmitReviewReq struct {
	Appropriateness bool   `json:"appropriateness"`
	Sincerity       bool   `json:"sincerity"`
	Reason          string `json:"reason" binding:"required,min=10,max=500"`
}

type ReviewSubmission struct {
	ID             int64     `json:"id"`
	BatchID        int64     `json:"batch_id"`
	Result         string    `json:"result"`
	LLMCheckStatus string    `json:"llm_check_status"`
	CompletedAt    time.Time `json:"completed_at"`
}
