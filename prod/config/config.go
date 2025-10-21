package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Broker   string
	Topic    string
	DataPath string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Ошибка загрузки .env файла, используйте конфиг по умолчанию")
	}
	return &Config{
		Broker:   os.Getenv("KAFKA_BROKER"),
		Topic:    os.Getenv("KAFKA_TOPIC"),
		DataPath: os.Getenv("PRODUCER_DATA_PATH"),
	}
}
