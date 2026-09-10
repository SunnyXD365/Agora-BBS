package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Agora-BBS/internal/model"
	"Agora-BBS/internal/pkg/response"
	"Agora-BBS/internal/service"
)

type LikeHandler struct {
	likeService *service.LikeService
}

func NewLikeHandler(s *service.LikeService) *LikeHandler {
	return &LikeHandler{likeService: s}
}

func (h *LikeHandler) Toggle(c *gin.Context) {
	userIDVal, _ := c.Get("current_user_id")
	userID := userIDVal.(int64)

	var req model.ToggleLikeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, err.Error())
		return
	}

	resp, err := h.likeService.ToggleLike(c.Request.Context(), userID, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40002, err.Error())
		return
	}

	response.Success(c, resp)
}
