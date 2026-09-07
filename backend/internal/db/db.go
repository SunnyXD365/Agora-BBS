package db

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// 读取指定目录的 sql 脚本并自动执行版本升迁
func RunMigrations(dbURL string, migrationsPath string) error {
	sourceURL := fmt.Sprintf("file://%s", migrationsPath)

	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		return fmt.Errorf("初始化 migrate 实例失败: %w", err)
	}
	defer m.Close()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("执行数据库 Migration 失败: %w", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("[Migration] 数据库已是最新版本，无需更新。")
	} else {
		log.Println("[Migration] 数据库结构自动迁移成功！")
	}

	return nil
}

// 初始化 GORM 数据库连接
func InitGORM(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}

	log.Println("[Database] PostgreSQL 连接成功。")
	return db, nil
}
