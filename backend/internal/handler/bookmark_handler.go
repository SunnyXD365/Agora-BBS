package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"agora-backend/internal/pkg/response"
	"agora-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type BookmarkHandler struct{ service *service.BookmarkService }

func NewBookmarkHandler(bookmarkService *service.BookmarkService) *BookmarkHandler {
	return &BookmarkHandler{service: bookmarkService}
}

func bookmarkTopicID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("topicId"), 10, 64)
	if err != nil || id < 1 {
		response.Error(c, http.StatusBadRequest, 40001, "invalid topic id")
		return 0, false
	}
	return id, true
}

func (h *BookmarkHandler) Create(c *gin.Context) {
	topicID, ok := bookmarkTopicID(c)
	if !ok {
		return
	}
	if err := h.service.Create(c.Request.Context(), c.GetInt64("userID"), topicID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Error(c, http.StatusNotFound, 40401, "topic not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, 50001, "failed to create bookmark")
		return
	}
	response.Success(c, gin.H{"bookmarked": true})
}

func (h *BookmarkHandler) Delete(c *gin.Context) {
	topicID, ok := bookmarkTopicID(c)
	if !ok {
		return
	}
	if err := h.service.Delete(c.Request.Context(), c.GetInt64("userID"), topicID); err != nil {
		response.Error(c, http.StatusInternalServerError, 50002, "failed to delete bookmark")
		return
	}
	response.Success(c, gin.H{"bookmarked": false})
}

func (h *BookmarkHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	data, err := h.service.List(c.Request.Context(), c.GetInt64("userID"), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50003, "failed to list bookmarks")
		return
	}
	response.Success(c, data)
}
