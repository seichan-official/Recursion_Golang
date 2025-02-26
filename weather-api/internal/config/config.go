package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// 環境変数をロードする関数
func LoadEnv() {
	err := godotenv.Load("internal/config/.env")
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}

// 環境変数を取得する関数
func GetEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
