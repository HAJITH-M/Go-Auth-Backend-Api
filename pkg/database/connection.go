package database

import (
	"fmt"
	"go-auth-backend-api/internal/config/env"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() error {

	env.Load()

	cfg := env.AppEnv

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
		cfg.DB_HOST,
		cfg.DB_USER,
		cfg.DB_PASSWORD,
		cfg.DB_NAME,
		cfg.DB_PORT,
	)
	// fmt.Println("DSN print:", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return err
	}

	// ✅ Get underlying sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// 🔥 CONNECTION POOL SETTINGS (VERY IMPORTANT)

	sqlDB.SetMaxOpenConns(25)                 // max DB connections
	sqlDB.SetMaxIdleConns(10)                 // idle connections
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // recycle connections
	sqlDB.SetConnMaxIdleTime(2 * time.Minute) // idle timeout

	DB = db

	return nil
}
