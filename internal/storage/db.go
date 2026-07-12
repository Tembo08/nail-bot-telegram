package storage

import (
	"log"
	"nail_bot/configs"
	"nail_bot/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *configs.Config) error {
	dsn := cfg.DBConnectionString
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=nail_bot port=5432 sslmode=disable"
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// Автоматическая миграция
	err = DB.AutoMigrate(
		&models.User{},
		&models.Booking{},
		&models.UserSession{},
	)
	if err != nil {
		return err
	}

	log.Println("✅ PostgreSQL подключена, миграция выполнена")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
