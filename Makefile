include .env
export

BINARY = cons-app

.PHONY: up down produce run migrate-up migrate-down help clean

up:
	docker-compose up -d

down:
	docker-compose down --volumes --remove-orphans
	@del cons-app 2>nul

produce:
	go run ./prod/main.go

run:
	@go build -o $(BINARY) ./cons/cmd
	./$(BINARY)

migrate-up:
	migrate -database ${POSTGRESQL_URL} -path db/migrations up

migrate-down:
	migrate -database ${POSTGRESQL_URL} -path db/migrations down

help:
	@echo "Доступные команды:"
	@echo "  make up             - запуск контейнеров"
	@echo "  make down           - останавить контейнеры"
	@echo "  make produce        - отправить сообщение в топик Kafka"
	@echo "  make run            - создать и запустить основное приложение"
	@echo "  make migrate-up     - накат миграций"
	@echo "  make migrate-down   - откат миграций"