package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewDB(Dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", Dsn)
	if err != nil {
		log.Fatal("Подключение к БД прервано:", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
