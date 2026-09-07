package router

import (
	"agora-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(topicHandler *handler.TopicHandler) *gin.Engine {
	r := gin.Default()

	// 健康检查
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	v1 := r.Group("/api/v1")
	{
		v1.POST("/topics", topicHandler.CreateTopic)
	}

	return r
}
