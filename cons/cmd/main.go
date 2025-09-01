package main

import (
	"l0/cons/internal/cache"
	"l0/cons/internal/config"
	"l0/cons/internal/consumer"
	"l0/cons/internal/handler"
	"l0/cons/internal/migrationsinit"
	"l0/cons/internal/repository"
	"l0/cons/internal/services"
	"l0/cons/pkg/db"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()

	db, err := db.NewDB(cfg)
	if err != nil {
		log.Println(err.Error())
	}

	migrationsinit.ApplyMigrations(db)

	cache := cache.NewCache()

	repo := repository.NewRepository(db, cache)
	service := services.NewService(repo)
	orderHandler := handler.NewHandlerTask(service)

	go consumer.Consume(repo, cache)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /order/{order_uid}", orderHandler.GetOrderHandler())
	mux.Handle("/", http.FileServer(http.Dir("cons/web/")))

	log.Fatal(http.ListenAndServe(":8082", mux))
}
