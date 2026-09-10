package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Agora-BBS/internal/model"
	"Agora-BBS/internal/pkg/response"
	"Agora-BBS/internal/service"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(s *service.PostService) *PostHandler {
	return &PostHandler{postService: s}
}

func (h *PostHandler) Create(c *gin.Context) {
	userIDVal, _ := c.Get("current_user_id")
	userID := userIDVal.(int64)

	var req model.CreatePostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, err.Error())
		return
	}

	post, err := h.postService.CreatePost(c.Request.Context(), userID, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40002, err.Error())
		return
	}

	response.Success(c, post)
}

func (h *PostHandler) List(c *gin.Context) {
	var req model.PostListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, err.Error())
		return
	}

	posts, err := h.postService.ListPosts(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50000, err.Error())
		return
	}

	response.Success(c, posts)
}
