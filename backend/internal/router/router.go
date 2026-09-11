package router

import (
	"github.com/gin-gonic/gin"

	"agora-backend/internal/config"
	"agora-backend/internal/handler"
	"agora-backend/internal/middleware"
)

type Handlers struct {
	UserHandler       *handler.UserHandler
	CategoryHandler   *handler.CategoryHandler
	TopicHandler      *handler.TopicHandler
	PostHandler       *handler.PostHandler
	LikeHandler       *handler.LikeHandler
	BookmarkHandler   *handler.BookmarkHandler
	GovernanceHandler *handler.GovernanceHandler
	FeedbackHandler   *handler.FeedbackHandler
	ReviewHandler     *handler.ReviewHandler
}

func SetupRouter(cfg *config.Config, h *Handlers) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.RequestID())

	v1 := r.Group("/api/v1")
	v1.Use(middleware.OptionalJWT(cfg.JWTSecret))
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
		v1.GET("/governance/policy", h.GovernanceHandler.Policy)
		v1.GET("/feedbacks/summary", h.FeedbackHandler.Summary)

		topics := v1.Group("/topics")
		{
			topics.GET("", h.TopicHandler.ListTopics)
			topics.GET("/:id", h.TopicHandler.GetTopicDetail)
			topics.GET("/:id/posts", h.PostHandler.ListPosts)
			topics.GET("/:id/clusters", h.FeedbackHandler.Clusters)
		}

		// 2. 受保护路由
		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			protected.GET("/users/me", h.UserHandler.GetProfile)
			protected.PUT("/users/me/onboarding", h.GovernanceHandler.SaveOnboarding)
			protected.POST("/reading-sessions", h.GovernanceHandler.StartReading)
			protected.PATCH("/reading-sessions/:id/heartbeat", h.GovernanceHandler.Heartbeat)
			protected.POST("/reading-sessions/:id/complete", h.GovernanceHandler.Complete)

			protected.POST("/topics", h.TopicHandler.CreateTopic)
			protected.PATCH("/topics/:id", h.TopicHandler.UpdateCooling)
			protected.DELETE("/topics/:id", h.TopicHandler.RecallCooling)
			protected.POST("/topics/:id/posts", h.PostHandler.CreatePost)
			protected.PATCH("/posts/:id", h.PostHandler.UpdateCooling)
			protected.DELETE("/posts/:id", h.PostHandler.RecallCooling)
			protected.GET("/bookmarks", h.BookmarkHandler.List)
			protected.PUT("/bookmarks/:topicId", h.BookmarkHandler.Create)
			protected.DELETE("/bookmarks/:topicId", h.BookmarkHandler.Delete)
			protected.POST("/feedbacks", h.FeedbackHandler.Upsert)
			protected.DELETE("/feedbacks/:id", h.FeedbackHandler.Withdraw)
			protected.GET("/reviews/tasks", h.ReviewHandler.List)
			protected.POST("/reviews/tasks/:id", h.ReviewHandler.Submit)

			protected.POST("/likes", h.LikeHandler.Like)
			protected.DELETE("/likes", h.LikeHandler.Unlike)
		}
	}

	return r
}
