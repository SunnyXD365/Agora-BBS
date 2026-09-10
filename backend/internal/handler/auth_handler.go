package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Agora-BBS/internal/model"
	"Agora-BBS/internal/pkg/response"
	"Agora-BBS/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register 处理用户注册请求 POST /api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "invalid request parameters: "+err.Error())
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40002, err.Error())
		return
	}

	response.Success(c, resp)
}

// Login 处理用户登录请求 POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "invalid request parameters: "+err.Error())
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, 40100, err.Error())
		return
	}

	response.Success(c, resp)
}

// GetMe 获取当前登录用户信息 GET /api/auth/me
func (h *AuthHandler) GetMe(c *gin.Context) {
	userIDVal, exists := c.Get("current_user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
		return
	}

	userID, ok := userIDVal.(int64)
	if !ok {
		response.Error(c, http.StatusInternalServerError, 50000, "invalid user_id type in context")
		return
	}

	user, err := h.authService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, 40400, err.Error())
		return
	}

	response.Success(c, user)
}
