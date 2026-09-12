package model

import "time"

type User struct {
	ID                  int64     `db:"id" json:"id"`
	Username            string    `db:"username" json:"username"`
	PasswordHash        string    `db:"password_hash" json:"-"` // 不在 JSON 中暴露
	Email               string    `db:"email" json:"email"`
	Avatar              string    `db:"avatar" json:"avatar"`
	Role                string    `db:"role" json:"role"`
	Status              string    `db:"status" json:"status"`
	TrustScore          int       `db:"trust_score" json:"-"`
	UnlockLevel         int       `db:"unlock_level" json:"unlock_level"`
	VerifiedReadSeconds int64     `db:"verified_read_seconds" json:"verified_read_seconds"`
	Capabilities        []string  `json:"capabilities" db:"-"`
	OnboardingStatement string    `json:"onboarding_statement" db:"onboarding_statement"`
	BackgroundTag       string    `json:"background_tag" db:"background_tag"`
	OnboardingStatus    string    `json:"onboarding_status" db:"onboarding_status"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time `db:"updated_at" json:"updated_at"`
}

type OnboardingReq struct {
	Statement     string `json:"statement" binding:"required,min=20,max=1000"`
	BackgroundTag string `json:"background_tag" binding:"required,min=2,max=64"`
}

// RegisterReq 注册请求体 DTO
type RegisterReq struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Email    string `json:"email" binding:"omitempty,email"`
}

// LoginReq 登录请求体 DTO
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type VerifyAdminEmailReq struct {
	ChallengeID string `json:"challenge_id" binding:"required,min=32,max=64"`
	Code        string `json:"code" binding:"required,len=6,numeric"`
}

type LoginResp struct {
	Token                       string `json:"token,omitempty"`
	User                        *User  `json:"user,omitempty"`
	RequiresEmailVerification   bool   `json:"requires_email_verification"`
	ChallengeID                 string `json:"challenge_id,omitempty"`
	MaskedEmail                 string `json:"masked_email,omitempty"`
	ExpiresInSeconds            int    `json:"expires_in_seconds,omitempty"`
	DevelopmentVerificationCode string `json:"development_verification_code,omitempty"`
}

// AuthResp 认证成功响应 DTO
type AuthResp struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}
