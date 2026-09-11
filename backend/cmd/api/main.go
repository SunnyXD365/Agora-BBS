package main

import (
	"context"
	"log"
	"time"

	"agora-backend/internal/cache"
	"agora-backend/internal/config"
	"agora-backend/internal/dao"
	"agora-backend/internal/db"
	"agora-backend/internal/handler"
	"agora-backend/internal/router"
	"agora-backend/internal/service"
	agoraworkflow "agora-backend/internal/workflow"
	"go.temporal.io/sdk/client"
)

func main() {
	// 1. 加载配置
	cfg := config.LoadConfig()
	log.Printf("[Init] Config loaded, API port: %s", cfg.Port)

	// 2. 初始化数据库连接与数据库迁移
	database, err := db.InitDB(cfg.DBDSN)
	if err != nil {
		log.Fatalf("[Error] Failed to initialize database: %v", err)
	}
	defer database.Close()
	log.Println("[Init] PostgreSQL connected successfully.")
	redisStore := cache.NewRedis(cfg.RedisAddr)
	defer redisStore.Close()
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 2*time.Second)
	if err := redisStore.Ping(pingCtx); err != nil {
		log.Printf("[Warn] Redis unavailable; requests will fall back to PostgreSQL: %v", err)
	} else {
		log.Println("[Init] Redis connected successfully.")
	}
	pingCancel()

	if err := db.RunMigrations(cfg.DBDSN); err != nil {
		log.Fatalf("[Error] Database migration failed: %v", err)
	}
	temporalClient, err := client.Dial(client.Options{HostPort: cfg.TemporalHost})
	if err != nil {
		log.Fatalf("[Error] Failed to connect to Temporal: %v", err)
	}
	defer temporalClient.Close()
	coolingStarter := agoraworkflow.NewTemporalStarter(temporalClient, cfg.TemporalTaskQueue)

	// 3. DAO 层初始化
	userDAO := dao.NewUserDAO(database)
	categoryDAO := dao.NewCategoryDAO(database)
	topicDAO := dao.NewTopicDAO(database)
	postDAO := dao.NewPostDAO(database)
	likeDAO := dao.NewLikeDAO(database)
	bookmarkDAO := dao.NewBookmarkDAO(database)
	governanceDAO := dao.NewGovernanceDAO(database)
	feedbackDAO := dao.NewFeedbackDAO(database)
	reviewDAO := dao.NewReviewDAO(database)
	adminDAO := dao.NewAdminDAO(database)

	// 4. Service 层初始化 (注入对应的 DAO 与配置项)
	userService := service.NewUserService(userDAO, cfg.JWTSecret, cfg.JWTExpireHours)
	categoryService := service.NewCategoryService(categoryDAO, redisStore)
	governanceService := service.NewGovernanceService(governanceDAO, userDAO, cfg, coolingStarter)
	topicService := service.NewTopicService(topicDAO, governanceService, coolingStarter, redisStore)
	postService := service.NewPostService(postDAO, governanceService, coolingStarter)
	likeService := service.NewLikeService(likeDAO)
	bookmarkService := service.NewBookmarkService(bookmarkDAO)
	feedbackService := service.NewFeedbackService(feedbackDAO, governanceService, coolingStarter)
	reviewService := service.NewReviewService(reviewDAO, governanceService, coolingStarter)
	adminService := service.NewAdminService(adminDAO, redisStore, coolingStarter, coolingStarter)

	// 5. Handler 层初始化 (注入对应的 Service)
	userHandler := handler.NewUserHandler(userService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	topicHandler := handler.NewTopicHandler(topicService)
	postHandler := handler.NewPostHandler(postService)
	likeHandler := handler.NewLikeHandler(likeService)
	bookmarkHandler := handler.NewBookmarkHandler(bookmarkService)
	governanceHandler := handler.NewGovernanceHandler(governanceService)
	feedbackHandler := handler.NewFeedbackHandler(feedbackService)
	reviewHandler := handler.NewReviewHandler(reviewService)
	adminHandler := handler.NewAdminHandler(adminService)

	// 6. 组装 Handlers 并传递给 SetupRouter
	handlers := &router.Handlers{
		UserHandler:       userHandler,
		CategoryHandler:   categoryHandler,
		TopicHandler:      topicHandler,
		PostHandler:       postHandler,
		LikeHandler:       likeHandler,
		BookmarkHandler:   bookmarkHandler,
		GovernanceHandler: governanceHandler,
		FeedbackHandler:   feedbackHandler,
		ReviewHandler:     reviewHandler,
		AdminHandler:      adminHandler,
	}

	r := router.SetupRouter(cfg, redisStore, database, handlers)

	// 7. 启动 HTTP 服务
	serverAddr := ":" + cfg.Port
	log.Printf("[Ready] agora-backend API server running on %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("[Fatal] Server forced to shutdown: %v", err)
	}
}
