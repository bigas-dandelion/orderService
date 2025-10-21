package main

import (
	"context"
	"log"

	"l0/prod/config"
	"l0/prod/generator"

	"github.com/segmentio/kafka-go"
)

func main() {
	cfg := config.LoadConfig()

	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Broker),
		Topic:    cfg.Topic,
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	data, err := generator.GenerateData()
	if err != nil {
		log.Fatalf("Ошибка в генерации заказа: %v", err)
	}

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Value: data,
	})

	if err != nil {
		log.Fatal("Не удалось отправить сообщение в топик:", err)
	}

	log.Println("Сообщение отправлено в топик 'orders'")
}
