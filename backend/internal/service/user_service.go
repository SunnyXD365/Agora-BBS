package service

import (
	"context"
	"errors"

	"agora-backend/internal/dao"
	"agora-backend/internal/model"
	"agora-backend/internal/pkg/hash"
	"agora-backend/internal/pkg/jwt"
)

type UserService struct {
	userDAO   *dao.UserDAO
	jwtSecret string
}

func NewUserService(userDAO *dao.UserDAO, jwtSecret string) *UserService {
	return &UserService{
		userDAO:   userDAO,
		jwtSecret: jwtSecret,
	}
}

func (s *UserService) Register(ctx context.Context, req *model.RegisterReq) (*model.AuthResp, error) {
	// 1. 业务逻辑校验：检查用户名是否已被占用
	existingUser, _ := s.userDAO.GetUserByUsername(ctx, req.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// 2. 调用 pkg/hash 加密密码
	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 3. 构造 Entity 落库
	user := &model.User{
		Username:     req.Username,
		PasswordHash: hashedPassword,
		Email:        req.Email,
	}
	if err := s.userDAO.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// 4. 颁发 JWT Token (有效期 72 小时)
	token, err := jwt.GenerateToken(user.ID, s.jwtSecret, 72)
	if err != nil {
		return nil, err
	}

	return &model.AuthResp{Token: token, User: user}, nil
}

func (s *UserService) Login(ctx context.Context, req *model.LoginReq) (*model.AuthResp, error) {
	// 1. 查找用户
	user, err := s.userDAO.GetUserByUsername(ctx, req.Username)
	if err != nil || user == nil {
		return nil, errors.New("invalid username or password")
	}

	// 2. 比对密码
	if !hash.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid username or password")
	}

	// 3. 生成 Token
	token, err := jwt.GenerateToken(user.ID, s.jwtSecret, 72)
	if err != nil {
		return nil, err
	}

	return &model.AuthResp{Token: token, User: user}, nil
}

func (s *UserService) GetProfile(ctx context.Context, userID int64) (*model.User, error) {
	return s.userDAO.GetUserByID(ctx, userID)
}