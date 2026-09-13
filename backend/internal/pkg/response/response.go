package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code      int    `json:"code"` // 业务自定义状态码 (0 表示成功)
	Msg       string `json:"msg"`  // 提示信息
	Data      any    `json:"data"` // 业务 Payload
	RequestID string `json:"request_id"`
}

// Success 成功响应 (HTTP 200)
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:      0,
		Msg:       "success",
		Data:      data,
		RequestID: c.GetString("requestID"),
	})
}

// Error 失败响应 (自定义 HTTP 状态码与业务码)
func Error(c *gin.Context, httpCode int, errCode int, msg string) {
	c.JSON(httpCode, Response{
		Code:      errCode,
		Msg:       msg,
		Data:      nil,
		RequestID: c.GetString("requestID"),
	})
}
