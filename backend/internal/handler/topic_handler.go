package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"agora-backend/internal/model"
	"agora-backend/internal/service"
	"agora-backend/internal/pkg/response"
)

type TopicHandler struct {
	topicService *service.TopicService
}

func NewTopicHandler(topicService *service.TopicService) *TopicHandler {
	return &TopicHandler{topicService: topicService}
}

func (h *TopicHandler) CreateTopic(c *gin.Context) {
	userID := c.GetInt64("userID")

	var req model.CreateTopicReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "invalid request parameters")
		return
	}

	topicID, err := h.topicService.CreateTopic(c.Request.Context(), userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50002, "failed to create topic")
		return
	}

	response.Success(c, gin.H{"topic_id": topicID})
}

func (h *TopicHandler) GetTopicDetail(c *gin.Context) {
	idStr := c.Param("id")
	topicID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40003, "invalid topic id")
		return
	}

	topic, err := h.topicService.GetTopicDetail(c.Request.Context(), topicID)
	if err != nil || topic == nil {
		response.Error(c, http.StatusNotFound, 40401, "topic not found")
		return
	}

	response.Success(c, topic)
}

func (h *TopicHandler) ListTopics(c *gin.Context) {
	categoryID, _ := strconv.ParseInt(c.Query("category_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	req := model.TopicListReq{
		CategoryID: categoryID,
		Page:       page,
		PageSize:   pageSize,
	}

	topics, err := h.topicService.ListTopics(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50003, "failed to fetch topics")
		return
	}

	response.Success(c, topics)
}