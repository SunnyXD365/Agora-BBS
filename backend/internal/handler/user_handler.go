package handler

import (
	"errors"
	"net/http"

	"agora-backend/internal/model"
	"agora-backend/internal/pkg/response"
	"agora-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req model.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "invalid request body: "+err.Error())
		return
	}

	resp, err := h.userService.Register(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40002, err.Error())
		return
	}

	response.Success(c, resp)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "invalid request body")
		return
	}

	resp, err := h.userService.Login(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrAdminEmailUnavailable) {
			response.Error(c, http.StatusServiceUnavailable, 50301, "administrator email service unavailable")
			return
		}
		response.Error(c, http.StatusUnauthorized, 40104, err.Error())
		return
	}

	response.Success(c, resp)
}

func (h *UserHandler) VerifyAdminEmail(c *gin.Context) {
	var req model.VerifyAdminEmailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40003, "invalid verification request")
		return
	}
	resp, err := h.userService.VerifyAdminEmail(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, 40105, "invalid or expired administrator verification code")
		return
	}
	response.Success(c, resp)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetInt64("userID") // 从中间件注入的上下文获取
	user, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "failed to fetch user profile")
		return
	}

	response.Success(c, user)
}
