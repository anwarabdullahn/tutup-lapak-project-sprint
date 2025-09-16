package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewGorm creates a new GORM database connection using PostgreSQL
func NewGorm(config *viper.Viper) *gorm.DB {
	// Try to get DATABASE_URL first, fallback to individual env vars
	dsn := config.GetString("DATABASE_URL")
	maxConnections := config.GetInt("DATABASE_MAXCONNECTIONS")

	// Fallback to individual database configuration if DATABASE_URL is not set
	if dsn == "" {
		username := config.GetString("DATABASE_USERNAME")
		password := config.GetString("DATABASE_PASSWORD")
		host := config.GetString("DATABASE_HOST")
		port := config.GetInt("DATABASE_PORTS")
		dbName := config.GetString("DATABASE_NAME")
		fmt.Println(dbName)
		if username == "" || host == "" || dbName == "" {
			log.Fatal("Database credentials are required")
		}

		if password == "" {
			dsn = fmt.Sprintf("host=%s user=%s dbname=%s port=%d sslmode=disable",
				host, username, dbName, port)
		} else {
			dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
				host, username, password, dbName, port)
		}
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get underlying sql.DB:", err)
	}

	// Configure connection pool
	if maxConnections > 0 {
		sqlDB.SetMaxOpenConns(maxConnections)
	} else {
		sqlDB.SetMaxOpenConns(10) // default
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Database connected successfully")
	return db
}
