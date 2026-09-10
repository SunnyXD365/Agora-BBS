package service

import (
	"context"
	"errors"

	"Agora-BBS/internal/config"
	"Agora-BBS/internal/dao"
	"Agora-BBS/internal/model"
	"Agora-BBS/internal/pkg/jwt"
	"Agora-BBS/internal/pkg/password"
)

type AuthService struct {
	userDAO *dao.UserDAO
	cfg     *config.Config
}

func NewAuthService(userDAO *dao.UserDAO, cfg *config.Config) *AuthService {
	return &AuthService{
		userDAO: userDAO,
		cfg:     cfg,
	}
}

// Register 用户注册逻辑
func (s *AuthService) Register(ctx context.Context, req *model.RegisterReq) (*model.AuthResp, error) {
	// 1. 检查用户名是否已存在
	existingUser, err := s.userDAO.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username already taken")
	}

	// 2. 密码加密
	pwdHash, err := password.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to encrypt password")
	}

	// 3. 构造实体并持久化
	newUser := &model.User{
		Username:     req.Username,
		PasswordHash: pwdHash,
		Email:        req.Email,
		Avatar:       "",
		Role:         "user",
		Status:       "active",
		TrustScore:   100,
		UnlockLevel:  1,
	}

	if err := s.userDAO.CreateUser(ctx, newUser); err != nil {
		return nil, err
	}

	// 4. 注册成功后直接生成 JWT Token
	token, err := jwt.GenerateToken(newUser.ID, newUser.Role, s.cfg.JWTSecret, s.cfg.JWTExpireHours)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &model.AuthResp{
		Token: token,
		User:  newUser,
	}, nil
}

// Login 用户登录逻辑
func (s *AuthService) Login(ctx context.Context, req *model.LoginReq) (*model.AuthResp, error) {
	// 1. 查找用户
	u, err := s.userDAO.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("invalid username or password")
	}

	// 2. 校验账号状态
	if u.Status == "banned" {
		return nil, errors.New("account has been banned")
	}

	// 3. 校验密码
	if !password.CheckPassword(req.Password, u.PasswordHash) {
		return nil, errors.New("invalid username or password")
	}

	// 4. 签发 Token
	token, err := jwt.GenerateToken(u.ID, u.Role, s.cfg.JWTSecret, s.cfg.JWTExpireHours)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &model.AuthResp{
		Token: token,
		User:  u,
	}, nil
}

// GetProfile 获取当前登录用户信息
func (s *AuthService) GetProfile(ctx context.Context, userID int64) (*model.User, error) {
	u, err := s.userDAO.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("user not found")
	}
	return u, nil
}
