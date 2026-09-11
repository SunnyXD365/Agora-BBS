package model

import (
	"encoding/json"
	"time"
)

type AdminDailyTrend struct {
	Date     string `json:"date"`
	Users    int    `json:"users"`
	Topics   int    `json:"topics"`
	Posts    int    `json:"posts"`
	Feedback int    `json:"feedback"`
}

type AdminOverview struct {
	UsersTotal        int64             `json:"users_total"`
	TopicsTotal       int64             `json:"topics_total"`
	PostsTotal        int64             `json:"posts_total"`
	ActiveUsers7Days  int64             `json:"active_users_7_days"`
	ContentStatus     map[string]int64  `json:"content_status"`
	TrustDistribution map[string]int64  `json:"trust_distribution"`
	ReviewTotal       int64             `json:"review_total"`
	ReviewCompleted   int64             `json:"review_completed"`
	ReviewExpired     int64             `json:"review_expired"`
	ReviewFair        int64             `json:"review_fair"`
	LLMCalls          int64             `json:"llm_calls"`
	LLMSuccess        int64             `json:"llm_success"`
	LLMAverageMS      float64           `json:"llm_average_ms"`
	LLMPromptTokens   int64             `json:"llm_prompt_tokens"`
	LLMOutputTokens   int64             `json:"llm_output_tokens"`
	Trend             []AdminDailyTrend `json:"trend"`
}

type AdminUser struct {
	ID                  int64     `json:"id"`
	Username            string    `json:"username"`
	Email               string    `json:"email"`
	Role                string    `json:"role"`
	Status              string    `json:"status"`
	UnlockLevel         int       `json:"unlock_level"`
	TrustScore          int       `json:"trust_score"`
	VerifiedReadSeconds int64     `json:"verified_read_seconds"`
	OnboardingStatus    string    `json:"onboarding_status"`
	CreatedAt           time.Time `json:"created_at"`
}

type AdminContent struct {
	ID         int64     `json:"id"`
	Type       string    `json:"type"`
	Title      string    `json:"title"`
	Excerpt    string    `json:"excerpt"`
	AuthorName string    `json:"author_name"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type AdminLLMJob struct {
	ID               int64           `json:"id"`
	JobType          string          `json:"job_type"`
	AggregateType    string          `json:"aggregate_type"`
	AggregateID      int64           `json:"aggregate_id"`
	Status           string          `json:"status"`
	Attempts         int             `json:"attempts"`
	Model            string          `json:"model"`
	Result           json.RawMessage `json:"result"`
	ErrorMessage     string          `json:"error_message"`
	PromptTokens     int             `json:"prompt_tokens"`
	CompletionTokens int             `json:"completion_tokens"`
	LatencyMS        int             `json:"latency_ms"`
	CreatedAt        time.Time       `json:"created_at"`
	CompletedAt      *time.Time      `json:"completed_at"`
}

type AdminTrustLog struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	Username      string    `json:"username"`
	EventType     string    `json:"event_type"`
	ScoreDelta    int       `json:"score_delta"`
	Reason        string    `json:"reason"`
	ReferenceType string    `json:"reference_type"`
	ReferenceID   *int64    `json:"reference_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type UpdateCategoryReq struct {
	Name           string `json:"name" binding:"required,min=2,max=64"`
	Slug           string `json:"slug" binding:"required,min=2,max=64"`
	Description    string `json:"description" binding:"max=500"`
	SortOrder      int    `json:"sort_order"`
	IsActive       bool   `json:"is_active"`
	RequiresReview bool   `json:"requires_review"`
}
