package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"agora-backend/internal/model"
	"agora-backend/internal/pkg/response"
	"agora-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type GovernanceHandler struct{ service *service.GovernanceService }

func NewGovernanceHandler(governanceService *service.GovernanceService) *GovernanceHandler {
	return &GovernanceHandler{service: governanceService}
}

func (h *GovernanceHandler) Policy(c *gin.Context) { response.Success(c, h.service.Policy()) }

func (h *GovernanceHandler) StartReading(c *gin.Context) {
	var req model.StartReadingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "invalid request body")
		return
	}
	data, err := h.service.StartReading(c.Request.Context(), c.GetInt64("userID"), req.TopicID)
	if errors.Is(err, sql.ErrNoRows) {
		response.Error(c, http.StatusNotFound, 40401, "topic not found")
		return
	}
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "failed to start reading")
		return
	}
	response.Success(c, data)
}

func (h *GovernanceHandler) Heartbeat(c *gin.Context) {
	var req model.ReadingHeartbeatReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "invalid heartbeat")
		return
	}
	data, err := h.service.Heartbeat(c.Request.Context(), c.GetInt64("userID"), c.Param("id"), &req)
	if errors.Is(err, sql.ErrNoRows) {
		response.Error(c, http.StatusNotFound, 40401, "reading session not found")
		return
	}
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50002, "failed to record heartbeat")
		return
	}
	response.Success(c, data)
}

func (h *GovernanceHandler) Complete(c *gin.Context) {
	data, err := h.service.CompleteReading(c.Request.Context(), c.GetInt64("userID"), c.Param("id"))
	if errors.Is(err, sql.ErrNoRows) {
		response.Error(c, http.StatusNotFound, 40401, "reading session not found")
		return
	}
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50003, "failed to complete reading")
		return
	}
	response.Success(c, data)
}

func (h *GovernanceHandler) SaveOnboarding(c *gin.Context) {
	var req model.OnboardingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "invalid onboarding profile")
		return
	}
	if err := h.service.SaveOnboarding(c.Request.Context(), c.GetInt64("userID"), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, 50004, "failed to save onboarding profile")
		return
	}
	response.Success(c, gin.H{"status": "pending_review"})
}
