package handler

import (
	"net/http"
	"strconv"

	"agora-backend/internal/model"
	"agora-backend/internal/pkg/response"
	"agora-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct{ service *service.AdminService }

func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{service: adminService}
}

func adminPage(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	return page, pageSize
}
func adminID(c *gin.Context) (int64, error) { return strconv.ParseInt(c.Param("id"), 10, 64) }

func (h *AdminHandler) Overview(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	data, err := h.service.Overview(c.Request.Context(), days)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50090, "failed to load administrator overview")
		return
	}
	response.Success(c, data)
}
func (h *AdminHandler) Users(c *gin.Context) {
	page, size := adminPage(c)
	data, err := h.service.Users(c.Request.Context(), page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50091, "failed to load users")
		return
	}
	response.Success(c, data)
}
func (h *AdminHandler) SetUserStatus(c *gin.Context) {
	id, err := adminID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40090, "invalid user id")
		return
	}
	var req struct {
		Status string `json:"status" binding:"required,oneof=active suspended"`
	}
	if err = c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40091, "invalid user status")
		return
	}
	if err = h.service.SetUserStatus(c.Request.Context(), c.GetInt64("userID"), id, req.Status); err != nil {
		response.Error(c, http.StatusConflict, 40990, err.Error())
		return
	}
	response.Success(c, gin.H{"status": req.Status})
}
func (h *AdminHandler) Contents(c *gin.Context) {
	page, size := adminPage(c)
	data, err := h.service.Contents(c.Request.Context(), c.DefaultQuery("type", "topic"), c.Query("status"), page, size)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40092, err.Error())
		return
	}
	response.Success(c, data)
}
func (h *AdminHandler) SetContentVisibility(c *gin.Context) {
	id, err := adminID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40093, "invalid content id")
		return
	}
	var req struct {
		Hidden *bool `json:"hidden" binding:"required"`
	}
	if err = c.ShouldBindJSON(&req); err != nil || req.Hidden == nil {
		response.Error(c, http.StatusBadRequest, 40094, "invalid visibility")
		return
	}
	if err = h.service.SetContentVisibility(c.Request.Context(), c.Param("type"), id, *req.Hidden); err != nil {
		response.Error(c, http.StatusConflict, 40991, err.Error())
		return
	}
	response.Success(c, gin.H{"hidden": *req.Hidden})
}
func (h *AdminHandler) LLMJobs(c *gin.Context) {
	page, size := adminPage(c)
	data, err := h.service.LLMJobs(c.Request.Context(), c.Query("status"), page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50092, "failed to load LLM jobs")
		return
	}
	response.Success(c, data)
}
func (h *AdminHandler) TrustLogs(c *gin.Context) {
	page, size := adminPage(c)
	userID, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	data, err := h.service.TrustLogs(c.Request.Context(), userID, page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50094, "failed to load trust logs")
		return
	}
	response.Success(c, data)
}
func (h *AdminHandler) RetryLLMJob(c *gin.Context) {
	id, err := adminID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40095, "invalid job id")
		return
	}
	if err = h.service.RetryLLMJob(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusConflict, 40992, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "queued"})
}
func (h *AdminHandler) Categories(c *gin.Context) {
	data, err := h.service.Categories(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50093, "failed to load categories")
		return
	}
	response.Success(c, data)
}
func (h *AdminHandler) CreateCategory(c *gin.Context) {
	var req model.UpdateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40096, "invalid category")
		return
	}
	data, err := h.service.CreateCategory(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusConflict, 40993, "category slug already exists")
		return
	}
	response.Success(c, data)
}
func (h *AdminHandler) UpdateCategory(c *gin.Context) {
	id, err := adminID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40097, "invalid category id")
		return
	}
	var req model.UpdateCategoryReq
	if err = c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40096, "invalid category")
		return
	}
	data, err := h.service.UpdateCategory(c.Request.Context(), id, &req)
	if err != nil {
		response.Error(c, http.StatusConflict, 40993, "category update failed")
		return
	}
	response.Success(c, data)
}
