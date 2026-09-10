package model

import "time"

type User struct {
	ID           int64     `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	PasswordHash string    `db:"password_hash" json:"-"` // 不在 JSON 中暴露
	Email        string    `db:"email" json:"email"`
	Avatar       string    `db:"avatar" json:"avatar"`
	Role         string    `db:"role" json:"role"`
	Status       string    `db:"status" json:"status"`
	TrustScore   int       `db:"trust_score" json:"trust_score"`
	UnlockLevel  int       `db:"unlock_level" json:"unlock_level"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
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

// AuthResp 认证成功响应 DTO
type AuthResp struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}
