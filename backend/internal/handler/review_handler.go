package handler

import (
	"net/http"
	"strconv"

	"agora-backend/internal/model"
	"agora-backend/internal/pkg/response"
	"agora-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ReviewHandler struct{ service *service.ReviewService }

func NewReviewHandler(s *service.ReviewService) *ReviewHandler { return &ReviewHandler{service: s} }
func (h *ReviewHandler) List(c *gin.Context) {
	data, err := h.service.List(c.Request.Context(), c.GetInt64("userID"))
	if err != nil {
		response.Error(c, http.StatusForbidden, 40301, err.Error())
		return
	}
	response.Success(c, data)
}
func (h *ReviewHandler) Submit(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "invalid review task id")
		return
	}
	var req model.SubmitReviewReq
	if err = c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40002, "invalid review submission")
		return
	}
	data, err := h.service.Submit(c.Request.Context(), c.GetInt64("userID"), id, &req)
	if err != nil {
		response.Error(c, http.StatusConflict, 40901, err.Error())
		return
	}
	response.Success(c, data)
}
