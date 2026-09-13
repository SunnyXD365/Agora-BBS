package handler

import (
	"net/http"

	"agora-backend/internal/model"
	"agora-backend/internal/pkg/response"
	"agora-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type LikeHandler struct {
	likeService *service.LikeService
}

func NewLikeHandler(likeService *service.LikeService) *LikeHandler {
	return &LikeHandler{likeService: likeService}
}

func (h *LikeHandler) Like(c *gin.Context) {
	userID := c.GetInt64("userID")
	var req model.ToggleLikeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "invalid request body")
		return
	}

	// 传递 &req 取指针
	if err := h.likeService.Like(c.Request.Context(), userID, &req); err != nil {
		response.Error(c, http.StatusBadRequest, 40002, err.Error())
		return
	}

	response.Success(c, gin.H{"status": "liked"})
}

func (h *LikeHandler) Unlike(c *gin.Context) {
	userID := c.GetInt64("userID")
	var req model.ToggleLikeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "invalid request body")
		return
	}

	// 传递 &req 取指针
	if err := h.likeService.Unlike(c.Request.Context(), userID, &req); err != nil {
		response.Error(c, http.StatusBadRequest, 40002, err.Error())
		return
	}

	response.Success(c, gin.H{"status": "unliked"})
}
