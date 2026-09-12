package router

import (
	"database/sql"
	"time"

	"agora-backend/internal/cache"
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
	AdminHandler      *handler.AdminHandler
	ContentHandler    *handler.ContentHandler
}

func SetupRouter(cfg *config.Config, store *cache.Store, database *sql.DB, h *Handlers) *gin.Engine {
	r := gin.Default()
	_ = r.SetTrustedProxies([]string{"127.0.0.1", "172.16.0.0/12"})
	r.Use(middleware.RequestID())
	r.Use(middleware.RateLimit(store, "api", 120000, time.Minute))

	v1 := r.Group("/api/v1")
	v1.Use(middleware.OptionalJWT(cfg.JWTSecret))
	{
		// 健康检查
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// 1. 公开路由
		auth := v1.Group("/auth")
		auth.Use(middleware.RateLimit(store, "auth", 30, time.Minute))
		{
			auth.POST("/register", h.UserHandler.Register)
			auth.POST("/login", h.UserHandler.Login)
			auth.POST("/admin/verify-email", h.UserHandler.VerifyAdminEmail)
		}

		v1.GET("/categories", h.CategoryHandler.ListCategories)
		v1.GET("/governance/policy", h.GovernanceHandler.Policy)
		v1.GET("/feedbacks/summary", h.FeedbackHandler.Summary)
		v1.GET("/search", h.ContentHandler.Search)

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
			protected.GET("/users/me/contents", h.ContentHandler.ListMine)
			protected.PUT("/users/me/onboarding", h.GovernanceHandler.SaveOnboarding)
			protected.POST("/reading-sessions", h.GovernanceHandler.StartReading)
			protected.PATCH("/reading-sessions/:id/heartbeat", h.GovernanceHandler.Heartbeat)
			protected.POST("/reading-sessions/:id/complete", h.GovernanceHandler.Complete)

			protected.POST("/topics", h.TopicHandler.CreateTopic)
			protected.POST("/topics/drafts", h.TopicHandler.CreateDraft)
			protected.PATCH("/topics/:id/draft", h.TopicHandler.UpdateDraft)
			protected.POST("/topics/:id/publish", h.TopicHandler.PublishDraft)
			protected.DELETE("/topics/:id/draft", h.TopicHandler.DeleteDraft)
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

		admin := v1.Group("/admin")
		admin.Use(middleware.JWTAuth(cfg.JWTSecret), middleware.AdminOnly(database))
		{
			admin.GET("/overview", h.AdminHandler.Overview)
			admin.GET("/users", h.AdminHandler.Users)
			admin.PATCH("/users/:id/status", h.AdminHandler.SetUserStatus)
			admin.GET("/contents", h.AdminHandler.Contents)
			admin.PATCH("/contents/:type/:id/visibility", h.AdminHandler.SetContentVisibility)
			admin.GET("/llm-jobs", h.AdminHandler.LLMJobs)
			admin.GET("/trust-logs", h.AdminHandler.TrustLogs)
			admin.POST("/llm-jobs/:id/retry", h.AdminHandler.RetryLLMJob)
			admin.GET("/categories", h.AdminHandler.Categories)
			admin.POST("/categories", h.AdminHandler.CreateCategory)
			admin.PATCH("/categories/:id", h.AdminHandler.UpdateCategory)
		}
	}

	return r
}
