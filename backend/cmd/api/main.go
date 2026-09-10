package main

import (
	"log"

	"Agora-BBS/internal/config"
	"Agora-BBS/internal/dao"
	"Agora-BBS/internal/db"
	"Agora-BBS/internal/handler"
	"Agora-BBS/internal/router"
	"Agora-BBS/internal/service"
)

func main() {
	cfg := config.LoadConfig()
	log.Printf("[Init] Config loaded, API port: %s", cfg.Port)

	database, err := db.InitDB(cfg.DBDSN)
	if err != nil {
		log.Fatalf("[Error] Failed to initialize database: %v", err)
	}
	defer database.Close()
	log.Println("[Init] PostgreSQL connected successfully.")

	if err := db.RunMigrations(cfg.DBDSN); err != nil {
		log.Fatalf("[Error] Database migration failed: %v", err)
	}
	// 1. DAO 初始化
	userDAO := dao.NewUserDAO(database)
	categoryDAO := dao.NewCategoryDAO(database)
	topicDAO := dao.NewTopicDAO(database)
	postDAO := dao.NewPostDAO(database)
	likeDAO := dao.NewLikeDAO(database)

	// 2. Service 初始化
	authService := service.NewAuthService(userDAO, cfg)
	categoryService := service.CategoryServiceFactory(categoryDAO)
	topicService := service.NewTopicService(topicDAO, categoryDAO)
	postService := service.NewPostService(postDAO, topicDAO)
	likeService := service.NewLikeService(likeDAO, topicDAO, postDAO)

	// 3. Handler 初始化
	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	topicHandler := handler.NewTopicHandler(topicService)
	postHandler := handler.NewPostHandler(postService)
	likeHandler := handler.NewLikeHandler(likeService)

	// 4. 路由与 HTTP 服务启动
	handlers := &router.Handlers{
		AuthHandler:     authHandler,
		CategoryHandler: categoryHandler,
		TopicHandler:    topicHandler,
		PostHandler:     postHandler,
		LikeHandler:     likeHandler,
	}

	r := router.SetupRouter(cfg, handlers)

	serverAddr := ":" + cfg.Port
	log.Printf("[Ready] Agora-BBS API server running on %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("[Fatal] Server forced to shutdown: %v", err)
	}
}
