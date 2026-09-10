package db

import (
	"errors"
	"fmt"
	"log"

	"Agora-BBS/migrations"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// RunMigrations 自动比对并执行增量 SQL 迁移脚本
func RunMigrations(dsn string) error {
	d, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("failed to init iofs migration source: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to execute migrations: %w", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("[Migration] Database schema is up to date.")
	} else {
		log.Println("[Migration] Database schema migrated successfully.")
	}

	return nil
}
