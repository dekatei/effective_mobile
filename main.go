package main

import (
	"fmt"
	"log"
	"net/http"

	"effective_mobile/base"
	"effective_mobile/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Подключение к базе данных
	db, err := base.CreateDB()
	if err != nil {
		log.Printf("Ошибка подключения к БД: %v", err)
		return
	}
	defer db.Close()

	// Создаём новый роутер
	r := chi.NewRouter()

	// Middleware (опционально, но полезно)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Роуты

	// Роуты подписок
	r.Post("/subscribe", handlers.HandlerAddSubscribe(db))           // Добавить подписку
	r.Put("/subscribe/{id}", handlers.HandlerUpdateSubscribe(db))    // Обновить подписку по ID
	r.Delete("/subscribe/{id}", handlers.HandlerDeleteSubscribe(db)) // Удалить подписку по ID

	r.Get("/subscribe/{id}", handlers.HandlerGetSubscribeByID(db))            // Получить одну подписку по ID
	r.Get("/subscribes/{user_id}", handlers.HandlerGetSubscribesByUserID(db)) // Получить все подписки пользователя
	// Роут стоимости
	r.Post("/cost", handlers.CostSummary(db)) // Суммарные траты, например GET /cost/abc123?start=2024-01&end=2024-07&service_name=Netflix

	// Запускаем сервер
	fmt.Println("Сервер запущен на http://localhost:8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
