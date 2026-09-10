package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"Agora-BBS/internal/model"
	"Agora-BBS/internal/pkg/response"
)

func TestTopicList_Validation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// 模拟帖子列表端点校验
	r.GET("/api/topics", func(c *gin.Context) {
		var req model.TopicListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			response.Error(c, http.StatusBadRequest, 40001, err.Error())
			return
		}
		response.Success(c, []model.Topic{})
	})

	// 测试默认参数请求
	req, _ := http.NewRequest(http.MethodGet, "/api/topics?page=1&page_size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"code":0`)
}

func TestCreateTopic_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/api/topics", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, 40101, "authorization header required")
			return
		}
	})

	// 测试未带 Token 提交发帖
	payload := model.CreateTopicReq{
		CategoryID: 1,
		Title:      "测试无鉴权发帖标题",
		Content:    "测试内容...",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/api/topics", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "40101")
}
