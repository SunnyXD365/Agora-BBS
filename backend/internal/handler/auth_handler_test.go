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

func TestRegisterValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 模拟注册参数校验逻辑
	r.POST("/api/auth/register", func(c *gin.Context) {
		var req model.RegisterReq
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, http.StatusBadRequest, 40001, "invalid request parameters")
			return
		}
		response.Success(c, gin.H{"status": "ok"})
	})

	// 1. 测试非法参数 (缺少 username/password)
	reqEmpty, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBufferString("{}"))
	reqEmpty.Header.Set("Content-Type", "application/json")
	wEmpty := httptest.NewRecorder()
	r.ServeHTTP(wEmpty, reqEmpty)

	assert.Equal(t, http.StatusBadRequest, wEmpty.Code)
	assert.Contains(t, wEmpty.Body.String(), "40001")

	// 2. 测试合法参数
	validReq := model.RegisterReq{
		Username: "testuser",
		Password: "password123",
	}
	body, _ := json.Marshal(validReq)
	reqValid, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
	reqValid.Header.Set("Content-Type", "application/json")
	wValid := httptest.NewRecorder()
	r.ServeHTTP(wValid, reqValid)

	assert.Equal(t, http.StatusOK, wValid.Code)
	assert.Contains(t, wValid.Body.String(), `"code":0`)
}
