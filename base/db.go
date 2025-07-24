package base

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// Создание базы данных 2ой шаг
func CreateDB() (db *sql.DB, err error) {
	connStr := os.Getenv("DB_CONN")
	if connStr == "" {
		connStr = "host=localhost port=8080 user=postgres password=postgres dbname=subscriptions sslmode=disable"
		log.Println("Используется стандартная строка подключения (без DB_CONN)")
	}

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Ошибка при подключении к БД: %v\n", err)
		return nil, err
	}

	// Проверка соединения
	if err := db.Ping(); err != nil {
		log.Printf("БД не доступна: %v\n", err)
		return nil, err
	}

	// Создание таблицы, если не существует
	if err := createSubscribesTable(db); err != nil {
		log.Printf("Ошибка при создании таблицы: %v\n", err)
		return nil, err
	}

	fmt.Println("База данных подключена и таблица subscriptions готова.")
	return db, nil
}

// Функция создания таблицы
func createSubscribesTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS subscribes (
		id SERIAL PRIMARY KEY,
		user_id UUID NOT NULL,
		service_name TEXT NOT NULL,
		price INTEGER NOT NULL CHECK (price >= 0),
		start_date DATE NOT NULL,
		expiration_date DATE
	);`

	_, err := db.Exec(query)
	return err
}
