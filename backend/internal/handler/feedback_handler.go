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

type FeedbackHandler struct{ service *service.FeedbackService }

func NewFeedbackHandler(s *service.FeedbackService) *FeedbackHandler {
	return &FeedbackHandler{service: s}
}
func (h *FeedbackHandler) Upsert(c *gin.Context) {
	var req model.UpsertFeedbackReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "invalid contextual feedback")
		return
	}
	data, err := h.service.Upsert(c.Request.Context(), c.GetInt64("userID"), &req)
	if errors.Is(err, service.ErrReplyLocked) {
		response.Error(c, http.StatusForbidden, 40301, err.Error())
		return
	}
	if errors.Is(err, dao.ErrSelfFeedback) {
		response.Error(c, http.StatusConflict, 40901, err.Error())
		return
	}
	if errors.Is(err, sql.ErrNoRows) {
		response.Error(c, http.StatusNotFound, 40401, "target not found")
		return
	}
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "failed to save feedback")
		return
	}
	response.Success(c, data)
}
func (h *FeedbackHandler) Withdraw(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "invalid feedback id")
		return
	}
	if err := h.service.Withdraw(c.Request.Context(), c.GetInt64("userID"), id); err != nil {
		response.Error(c, http.StatusConflict, 40901, "feedback is not withdrawable")
		return
	}
	response.Success(c, gin.H{"status": "withdrawn"})
}
func (h *FeedbackHandler) Summary(c *gin.Context) {
	targetType := c.Query("target_type")
	targetID, err := strconv.ParseInt(c.Query("target_id"), 10, 64)
	if err != nil || (targetType != "topic" && targetType != "post") {
		response.Error(c, http.StatusBadRequest, 40001, "invalid target")
		return
	}
	data, err := h.service.Summary(c.Request.Context(), targetType, targetID, c.GetInt64("userID"))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50002, "failed to load feedback summary")
		return
	}
	response.Success(c, data)
}
func (h *FeedbackHandler) Clusters(c *gin.Context) {
	topicID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "invalid topic id")
		return
	}
	data, err := h.service.Clusters(c.Request.Context(), topicID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50003, "failed to load clusters")
		return
	}
	response.Success(c, data)
}
