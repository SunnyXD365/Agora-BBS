package handler

import (
	"net/http"

	"agora-backend/internal/pkg/response"
	"agora-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type TopicHandler struct {
	topicService *service.TopicService
}

func NewTopicHandler(topicService *service.TopicService) *TopicHandler {
	return &TopicHandler{topicService: topicService}
}

func (h *TopicHandler) CreateTopic(c *gin.Context) {
	var req service.CreateTopicReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "请求参数无效: "+err.Error())
		return
	}

	// Mock 用户身份，后续接入 JWT 中间件提取
	var currentUserID int64 = 1

	topic, err := h.topicService.CreateTopic(c.Request.Context(), currentUserID, &req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "发帖失败: "+err.Error())
		return
	}

	response.Success(c, topic)
}
