package configs

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBotToken   string
	DBConnectionString string
	AdminID            int64
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env файл не найден")
	}

	adminID, _ := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)

	return &Config{
		TelegramBotToken:   os.Getenv("TELEGRAM_BOT_TOKEN"),
		DBConnectionString: os.Getenv("DB_CONNECTION_STRING"),
		AdminID:            adminID,
	}
}

func getEnvAsInt64(key string, defaultVal int64) int64 {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	var result int64
	_, err := fmt.Sscanf(val, "%d", &result)
	if err != nil {
		return defaultVal
	}
	return result
}
