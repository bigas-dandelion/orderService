package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"l0/cons/internal/cache"
	"l0/cons/internal/models"
	"l0/cons/internal/repository"
	"log"

	"github.com/segmentio/kafka-go"
)

const (
	kafkaTopic  = "orders"
	kafkaBroker = "localhost:29092"
	kafkaGroup  = "order-service-group"
)

func Consume(repo *repository.Repository, orderCache *cache.Cache) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{kafkaBroker},
		Topic:     kafkaTopic,
		GroupID:   kafkaGroup,
		Partition: 0,
		MinBytes:  10e3,
		MaxBytes:  10e6,
	})

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Ошибка при получении сообщения из Kafka: %v", err)
			continue
		}

		fmt.Printf("Полученное сообщение: %s\n", string(msg.Value))

		saveOrderInCacheAndDB(msg.Value, repo, orderCache)

		err = reader.CommitMessages(context.Background(), msg)
		if err != nil {
			log.Printf("Ошибка при подтверждении смещения сообщения: %v", err)
		}
	}
}

func saveOrderInCacheAndDB(message []byte, repo *repository.Repository, cache *cache.Cache) {
	var order models.Order

	err := json.Unmarshal(message, &order)
	if err != nil {
		log.Printf("Ошибка десериализации сообщения Kafka: %v", err)
		return
	}

	err = repo.StoreOrder(&order)
	if err != nil {
		log.Printf("Ошибка сохранения заказа %s в бд: %v", order.OrderUID, err)
		return
	}

	cache.Set(order.OrderUID, &order)
	log.Printf("Заказ %s обработан и сохранен в бд и кэше.", order.OrderUID)
}
