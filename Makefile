.PHONY: up down produce run help

up:
	docker-compose up -d

down:
	docker-compose down --volumes --remove-orphans

produce:
	go run ./prod/main.go

run:
	go run ./cons/cmd/main.go

help:
	@echo "Доступные команды:"
	@echo "  make up             - запустить контейнеры"
	@echo "  make down           - остановить контейнеры"
	@echo "  make produce        - отправить задачу в топик"
	@echo "  make run   	 	 - запуск основной программы"