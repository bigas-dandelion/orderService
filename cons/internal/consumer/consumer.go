package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"l0/cons/internal/cache"
	"l0/cons/internal/models"
	"l0/cons/internal/repository"
	"l0/cons/internal/validation"

	"github.com/segmentio/kafka-go"
)

func Consume(ctx context.Context, broker, topic, group string, repo *repository.Repository, orderCache *cache.Cache) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  group,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Println("Консьюмер кафки остановлен")
				return
			}
			log.Printf("Ошибка при получении сообщения из Kafka: %v", err)
			continue
		}

		fmt.Printf("\nПолученное сообщение: %s\n\n", msg.Value)
		saveOrderInCacheAndDB(msg.Value, repo, orderCache)

		if err = reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("Ошибка при подтверждении смещения: %v", err)
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

	if err = validation.Validate(&order); err != nil {
		log.Printf("Ошибка валидации заказа: %v", err)
		return
	}

	err = repo.StoreOrder(&order)
	if err != nil {
		log.Printf("Ошибка сохранения заказа %s в бд: %v", order.OrderUID, err)
		return
	}

	cache.Put(order.OrderUID, &order)
	log.Printf("Заказ %s обработан и сохранен в бд и кэше.", order.OrderUID)
}
