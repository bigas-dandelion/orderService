package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Broker   string
	Topic    string
}

func LoadConfig() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Ошибка при загрузке файла .env, используйте конфиг по умолчанию")
	}

	return &Config{
		Broker:   getEnvOrDefault("KAFKA_BROKER", "localhost:29092"),
		Topic:    getEnvOrDefault("KAFKA_TOPIC", "orders"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
