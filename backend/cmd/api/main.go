package main

import (
	"log"
	"os"

	"agora-backend/internal/dao"
	"agora-backend/internal/db"
	"agora-backend/internal/handler"
	"agora-backend/internal/router"
	"agora-backend/internal/service"

	"go.temporal.io/sdk/client"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://agora_user:agora_password@agora-postgres:5432/agora_db?sslmode=disable"
	}

	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "migrations"
	}

	temporalHost := os.Getenv("TEMPORAL_HOST")
	if temporalHost == "" {
		temporalHost = "agora-temporal:7233"
	}

	// 执行数据库 Migration
	if err := db.RunMigrations(dbURL, migrationsPath); err != nil {
		log.Fatalf("Migration 失败: %v", err)
	}

	// 初始化 GORM
	gormDB, err := db.InitGORM(dbURL)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 初始化 Temporal Client
	temporalClient, err := client.Dial(client.Options{HostPort: temporalHost})
	if err != nil {
		log.Fatalf("连接 Temporal 失败: %v", err)
	}
	defer temporalClient.Close()

	// 依赖注入 (DAO -> Service -> Handler)
	topicDAO := dao.NewTopicDAO(gormDB)
	topicService := service.NewTopicService(topicDAO, temporalClient)
	topicHandler := handler.NewTopicHandler(topicService)

	// 启动 Gin HTTP 服务
	r := router.SetupRouter(topicHandler)
	log.Println("API 服务器正在监听 :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
