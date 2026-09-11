package router

import (
	"github.com/gin-gonic/gin"

	"agora-backend/internal/config"
	"agora-backend/internal/handler"
	"agora-backend/internal/middleware"
)

type Handlers struct {
	UserHandler     *handler.UserHandler
	CategoryHandler *handler.CategoryHandler
	TopicHandler    *handler.TopicHandler
	PostHandler     *handler.PostHandler
	LikeHandler     *handler.LikeHandler
	BookmarkHandler *handler.BookmarkHandler
}

func SetupRouter(cfg *config.Config, h *Handlers) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.RequestID())

	v1 := r.Group("/api/v1")
	{
		// 健康检查
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// 1. 公开路由
		auth := v1.Group("/auth")
		{
			auth.POST("/register", h.UserHandler.Register)
			auth.POST("/login", h.UserHandler.Login)
		}

		v1.GET("/categories", h.CategoryHandler.ListCategories)

		topics := v1.Group("/topics")
		{
			topics.GET("", h.TopicHandler.ListTopics)
			topics.GET("/:id", h.TopicHandler.GetTopicDetail)
			topics.GET("/:id/posts", h.PostHandler.ListPosts)
		}

		// 2. 受保护路由
		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			protected.GET("/users/me", h.UserHandler.GetProfile)

			protected.POST("/topics", h.TopicHandler.CreateTopic)
			protected.POST("/topics/:id/posts", h.PostHandler.CreatePost)
			protected.GET("/bookmarks", h.BookmarkHandler.List)
			protected.PUT("/bookmarks/:topicId", h.BookmarkHandler.Create)
			protected.DELETE("/bookmarks/:topicId", h.BookmarkHandler.Delete)

			protected.POST("/likes", h.LikeHandler.Like)
			protected.DELETE("/likes", h.LikeHandler.Unlike)
		}
	}

	return r
}
