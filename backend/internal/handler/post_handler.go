package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"agora-backend/internal/dao"
	"agora-backend/internal/model"
	"agora-backend/internal/pkg/response"
	"agora-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	topicIDStr := c.Param("id")
	topicID, err := strconv.ParseInt(topicIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "invalid topic id")
		return
	}

	var req model.CreatePostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40002, "invalid request body")
		return
	}

	req.TopicID = topicID
	userID := c.GetInt64("userID")

	// 这里对齐 Service 的 (ctx, userID, &req) 签名
	postID, err := h.postService.CreatePost(c.Request.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Error(c, http.StatusNotFound, 40401, "topic not found")
			return
		}
		if errors.Is(err, dao.ErrInvalidParent) {
			response.Error(c, http.StatusBadRequest, 40003, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, 50002, err.Error())
		return
	}

	response.Success(c, gin.H{"post_id": postID})
}

func (h *PostHandler) ListPosts(c *gin.Context) {
	topicID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "invalid topic id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	pageData, err := h.postService.ListPosts(c.Request.Context(), topicID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "failed to fetch posts")
		return
	}

	response.Success(c, pageData)
}
