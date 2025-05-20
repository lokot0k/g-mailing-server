package database

import (
	"fmt"

	"github.com/lokot0k/mservice/config"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}

func Migrate(cfg *config.Config) error {
	goose.SetVerbose(true)
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
	sqlDB, err := goose.OpenDBWithDriver("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("goose open: %w", err)
	}
	defer sqlDB.Close()

	if err := goose.Up(sqlDB, "./db/migrations"); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	return nil
}
