package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewGorm initializes a Postgres connection using DATABASE_URL (preferred)
// and applies sane pooling defaults.
func NewGorm(cfg *viper.Viper) *gorm.DB {
	dsn := cfg.GetString("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required for product-service")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("database ping failed: %v", err)
	}
	return db
}
