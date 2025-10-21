package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Db       DbConfig
	KafkaCfg KafkaConfig
	CacheCfg CacheConfig
	HTTPCfg  HTTP
}

type KafkaConfig struct {
	Broker string
	Topic  string
	Group  string
}

type CacheConfig struct {
	Size int
}

type HTTP struct {
	Port string
}

type DbConfig struct {
	Dsn string
}

func LoadConfig() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Ошибка при загрузке файла .env, используйте конфиг по умолчанию")
	}

	return &Config{
		Db: DbConfig{
			Dsn: getEnvOrDefault("DSN", "host=localhost user=order_user password=order_pass dbname=orders_db port=5434 sslmode=disable"),
		},

		KafkaCfg: KafkaConfig{
			Broker: getEnvOrDefault("KAFKA_BROKER", "localhost:29092"),
			Topic:  getEnvOrDefault("KAFKA_TOPIC", "orders"),
			Group:  getEnvOrDefault("KAFKA_GROUP", "order-service-group"),
		},

		CacheCfg: CacheConfig{
			Size: getEnvIntOrDefault("CACHE_SIZE", 64),
		},

		HTTPCfg: HTTP{
			Port: getEnvOrDefault("HTTP_PORT", "8082"),
		},
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if valueStr := os.Getenv(key); valueStr != "" {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
		log.Printf("Недопустимое целое значение для %s: %s, испоьзуем по умолчанию %d", key, valueStr, defaultValue)
	}
	return defaultValue
}
