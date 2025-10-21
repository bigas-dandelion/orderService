package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"l0/cons/internal/cache"
	"l0/cons/internal/config"
	"l0/cons/internal/consumer"
	"l0/cons/internal/handler"
	"l0/cons/internal/repository"
	"l0/cons/internal/services"
	"l0/cons/pkg/db"
)

func main() {
	cfg := config.LoadConfig()

	dbConn, err := db.NewDB(cfg.Db.Dsn)
	if err != nil {
		log.Fatal("не удалось подключиться к базе данных:", err)
	}
	defer dbConn.Close()

	cacheStore := cache.NewCache(cfg.CacheCfg.Size)
	repo := repository.NewRepository(dbConn, cacheStore)
	service := services.NewService(repo)
	orderHandler := handler.NewHandlerTask(service)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go consumer.Consume(ctx, cfg.KafkaCfg.Broker, cfg.KafkaCfg.Topic, cfg.KafkaCfg.Group, repo, cacheStore)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /order/{order_uid}", orderHandler.GetOrderHandler())
	mux.Handle("/", http.FileServer(http.Dir("cons/web/")))

	srv := &http.Server{
		Addr:    ":" + cfg.HTTPCfg.Port,
		Handler: mux,
	}

	go func() {
		log.Printf("Сервер запущен на порту %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Ошибка сервера: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Безопасное завершение работы...")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Ошибка завершения сервера: %v", err)
	} else {
		log.Println("Сервер остановлен")
	}

	log.Println("Приложение остановлено")
}
