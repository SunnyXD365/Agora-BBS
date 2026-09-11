package main

import (
	"log"

	"agora-backend/internal/config"
	"agora-backend/internal/dao"
	"agora-backend/internal/db"
	"agora-backend/internal/handler"
	"agora-backend/internal/router"
	"agora-backend/internal/service"
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

	if err := db.RunMigrations(cfg.DBDSN); err != nil {
		log.Fatalf("[Error] Database migration failed: %v", err)
	}

	// 3. DAO 层初始化
	userDAO := dao.NewUserDAO(database)
	categoryDAO := dao.NewCategoryDAO(database)
	topicDAO := dao.NewTopicDAO(database)
	postDAO := dao.NewPostDAO(database)
	likeDAO := dao.NewLikeDAO(database)

	// 4. Service 层初始化 (注入对应的 DAO 与配置项)
	userService := service.NewUserService(userDAO, cfg.JWTSecret)
	categoryService := service.NewCategoryService(categoryDAO)
	topicService := service.NewTopicService(topicDAO)
	postService := service.NewPostService(postDAO)
	likeService := service.NewLikeService(likeDAO)

	// 5. Handler 层初始化 (注入对应的 Service)
	userHandler := handler.NewUserHandler(userService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	topicHandler := handler.NewTopicHandler(topicService)
	postHandler := handler.NewPostHandler(postService)
	likeHandler := handler.NewLikeHandler(likeService)

	// 6. 组装 Handlers 并传递给 SetupRouter
	handlers := &router.Handlers{
		UserHandler:     userHandler,
		CategoryHandler: categoryHandler,
		TopicHandler:    topicHandler,
		PostHandler:     postHandler,
		LikeHandler:     likeHandler,
	}

	r := router.SetupRouter(cfg, handlers)

	// 7. 启动 HTTP 服务
	serverAddr := ":" + cfg.Port
	log.Printf("[Ready] agora-backend API server running on %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("[Fatal] Server forced to shutdown: %v", err)
	}
}