package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Agora-BBS/internal/model"
	"Agora-BBS/internal/pkg/response"
	"Agora-BBS/internal/service"
)

type TopicHandler struct {
	topicService *service.TopicService
}

func NewTopicHandler(s *service.TopicService) *TopicHandler {
	return &TopicHandler{topicService: s}
}

func (h *TopicHandler) List(c *gin.Context) {
	var req model.TopicListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, err.Error())
		return
	}
	topics, err := h.topicService.ListTopics(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50000, err.Error())
		return
	}
	response.Success(c, topics)
}

func (h *TopicHandler) Create(c *gin.Context) {
	userIDVal, _ := c.Get("current_user_id")
	userID := userIDVal.(int64)

	var req model.CreateTopicReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, err.Error())
		return
	}

	topic, err := h.topicService.CreateTopic(c.Request.Context(), userID, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40002, err.Error())
		return
	}
	response.Success(c, topic)
}

func (h *TopicHandler) GetDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "invalid topic id")
		return
	}

	topic, err := h.topicService.GetTopicDetail(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, 40400, err.Error())
		return
	}
	response.Success(c, topic)
}
