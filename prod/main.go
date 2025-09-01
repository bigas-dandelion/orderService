package main

import (
	"context"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

const (
	kafkaTopic = "orders"
	kafkaBroker = "localhost:29092"
	dataPass = "prod/model.json"
)

func main() {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    kafkaTopic,
		Balancer: &kafka.LeastBytes{},
	}

	defer writer.Close()

	data, err := os.ReadFile(dataPass)
	if err != nil {
		log.Fatal("Не удалось прочитать файл:", err)
	}

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Value: data,
	})

	if err != nil {
		log.Fatal("Не удалось отправить сообщение в топик:", err)
	}

	log.Println("Сообщение отправлено в топик 'orders'")
}