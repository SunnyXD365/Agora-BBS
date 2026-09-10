package router

import (
	"github.com/gin-gonic/gin"

	"Agora-BBS/internal/config"
	"Agora-BBS/internal/handler"
	"Agora-BBS/internal/middleware"
)

type Handlers struct {
	AuthHandler     *handler.AuthHandler
	CategoryHandler *handler.CategoryHandler
	TopicHandler    *handler.TopicHandler
	PostHandler     *handler.PostHandler
	LikeHandler     *handler.LikeHandler
}

func SetupRouter(cfg *config.Config, handlers *Handlers) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong", "status": "healthy"})
		})

		// 认证公开接口
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", handlers.AuthHandler.Register)
			authGroup.POST("/login", handlers.AuthHandler.Login)
		}

		// 公开只读接口
		api.GET("/categories", handlers.CategoryHandler.List)
		api.GET("/topics", handlers.TopicHandler.List)
		api.GET("/topics/:id", handlers.TopicHandler.GetDetail)

		// 受保护接口 (JWT 认证)
		protected := api.Group("")
		protected.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			protected.GET("/auth/me", handlers.AuthHandler.GetMe)
			protected.POST("/categories", handlers.CategoryHandler.Create)
			protected.POST("/topics", handlers.TopicHandler.Create)
			protected.POST("/posts", handlers.PostHandler.Create)   // 发布回复
			protected.POST("/likes/toggle", handlers.LikeHandler.Toggle) // 点赞/取消点赞
		}
	}

	return r
}
