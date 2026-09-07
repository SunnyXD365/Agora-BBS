package main

import (
	"log"
	"os"

	"agora-backend/internal/db"
)

func main() {
	// 从环境变量获取数据库连接串
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Postgres DSN 格式
		dbURL = "postgres://agora_user:agora_password@agora-postgres:5432/agora_db?sslmode=disable"
	}

	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "migrations" // 相对路径，指向 backend/migrations
	}

	// 自动执行数据库 Migration 迁移
	log.Println("正在检查并同步数据库结构...")
	if err := db.RunMigrations(dbURL, migrationsPath); err != nil {
		log.Fatalf("数据库 Migration 失败，停止启动: %v", err)
	}

	// 初始化 GORM 数据库实例
	gormDB, err := db.InitGORM(dbURL)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	_ = gormDB // 后续将 gormDB 传给 Service/DAO 层使用

	log.Println("服务器启动成功，开始监听请求...")
	// TODO: 启动 Gin 或 HTTP Server
}
