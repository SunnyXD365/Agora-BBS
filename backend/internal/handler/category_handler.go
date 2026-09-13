package handler

import (
	"net/http"

	"agora-backend/internal/pkg/response"
	"agora-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

func (h *CategoryHandler) ListCategories(c *gin.Context) {
	categories, err := h.categoryService.ListCategories(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "failed to fetch categories")
		return
	}
	response.Success(c, categories)
}
