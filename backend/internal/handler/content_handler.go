package handler

import (
	"errors"
	"net/http"
	"strconv"

	"agora-backend/internal/pkg/response"
	"agora-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ContentHandler struct {
	service *service.ContentService
}

func NewContentHandler(contentService *service.ContentService) *ContentHandler {
	return &ContentHandler{service: contentService}
}

func (h *ContentHandler) Search(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	data, err := h.service.Search(c.Request.Context(), c.Query("q"), c.Query("type"), page, pageSize)
	if errors.Is(err, service.ErrInvalidContentQuery) {
		response.Error(c, http.StatusBadRequest, 40001, "search query must contain 2 to 100 characters")
		return
	}
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50003, "failed to search content")
		return
	}
	response.Success(c, data)
}

func (h *ContentHandler) ListMine(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	data, err := h.service.ListMine(c.Request.Context(), c.GetInt64("userID"), c.Query("type"), c.Query("status"), page, pageSize)
	if errors.Is(err, service.ErrInvalidContentQuery) {
		response.Error(c, http.StatusBadRequest, 40001, "invalid content filter")
		return
	}
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50003, "failed to fetch user content")
		return
	}
	response.Success(c, data)
}
