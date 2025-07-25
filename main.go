package main

import (
	"fmt"
	"log"
	"net/http"

	"effective_mobile/base"
	"effective_mobile/handlers"

	_ "effective_mobile/docs" // docs генерируется автоматически

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Effective Mobile Subscription API
// @version         1.0
// @description     API для управления подписками пользователя и подсчёта стоимости.

// @contact.name   API Support
// @contact.email  support@example.com

// @host      localhost:8080
// @BasePath  /

// @schemes http
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
	// Подключение к базе данных
	db, err := base.CreateDB()
	if err != nil {
		log.Printf("Ошибка подключения к БД: %v", err)
		return
	}
	defer db.Close()

	// Создаём новый роутер
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// @Summary      Добавить подписку
	// @Description  Создает новую подписку пользователя
	// @Tags         subscriptions
	// @Accept       json
	// @Produce      json
	// @Param        subscription  body      base.Subscription  true  "Подписка"
	// @Success      200           {object}  map[string]interface{}
	// @Failure      400           {object}  map[string]string
	// @Failure      500           {object}  map[string]string
	// @Router       /subscription [post]
	r.Post("/subscription", handlers.HandlerAddSubscription(db))

	// @Summary      Обновить подписку по ID
	// @Description  Обновляет подписку с указанным ID
	// @Tags         subscriptions
	// @Accept       json
	// @Produce      json
	// @Param        id            path      string             true  "ID подписки"
	// @Param        subscription  body      base.Subscription  true  "Подписка"
	// @Success      200           {object}  map[string]string
	// @Failure      400           {object}  map[string]string
	// @Failure      500           {object}  map[string]string
	// @Router       /subscription/{id} [put]
	r.Put("/subscription/{id}", handlers.HandlerUpdateSubscription(db))

	// @Summary      Удалить подписку по ID
	// @Description  Удаляет подписку с указанным ID
	// @Tags         subscriptions
	// @Produce      json
	// @Param        id  path  string  true  "ID подписки"
	// @Success      200 {object} map[string]interface{}
	// @Failure      400 {object} map[string]string
	// @Failure      500 {object} map[string]string
	// @Router       /subscription/{id} [delete]
	r.Delete("/subscription/{id}", handlers.HandlerDeleteSubscription(db))

	// @Summary      Получить подписку по ID
	// @Description  Возвращает подписку по ID
	// @Tags         subscriptions
	// @Produce      json
	// @Param        id  path  string  true  "ID подписки"
	// @Success      200 {object} base.Subscription
	// @Failure      400 {object} map[string]string
	// @Failure      404 {object} map[string]string
	// @Router       /subscription/{id} [get]
	r.Get("/subscription/{id}", handlers.HandlerGetSubscriptionByID(db))

	// @Summary      Получить подписки пользователя
	// @Description  Возвращает все подписки пользователя по user_id
	// @Tags         subscriptions
	// @Produce      json
	// @Param        user_id  path  string  true  "ID пользователя"
	// @Success      200      {array} base.Subscription
	// @Failure      400      {object} map[string]string
	// @Router       /subscriptions/{user_id} [get]
	r.Get("/subscriptions/{user_id}", handlers.HandlerGetSubscriptionsByUserID(db)) // Получить все подписки пользователя

	// @Summary      Получить суммарную стоимость подписок
	// @Description  Возвращает сумму затрат пользователя за период (с фильтром по сервису)
	// @Tags         cost
	// @Accept       json
	// @Produce      json
	// @Param        body  body  handlers.CostRequest  true  "Параметры фильтра"
	// @Success      200   {object} handlers.CostResponse
	// @Failure      400   {object} map[string]string
	// @Failure      500   {object} map[string]string
	// @Router       /cost [post]
	r.Get("/cost/{user_id}", handlers.CostSummary(db)) // Суммарные траты, например GET /cost/abc123?start=2024-01&end=2024-07&service_name=Netflix

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	// Запускаем сервер
	fmt.Println("Сервер запущен на http://localhost:8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
