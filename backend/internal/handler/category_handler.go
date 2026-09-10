package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Agora-BBS/internal/model"
	"Agora-BBS/internal/pkg/response"
	"Agora-BBS/internal/service"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(s *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: s}
}

func (h *CategoryHandler) List(c *gin.Context) {
	categories, err := h.categoryService.ListCategories(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50000, err.Error())
		return
	}
	response.Success(c, categories)
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req model.CreateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, err.Error())
		return
	}
	cat, err := h.categoryService.CreateCategory(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40002, err.Error())
		return
	}
	response.Success(c, cat)
}
