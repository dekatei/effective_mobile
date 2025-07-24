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
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, name)

	db, err = sql.Open("postgres", dsn)
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
		end_date DATE
	);`

	_, err := db.Exec(query)
	return err
}
