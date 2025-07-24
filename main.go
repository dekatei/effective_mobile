package main

import (
	"fmt"
	"log"
	"net/http"

	"effective_mobile/base"
	"effective_mobile/handlers"

	"github.com/go-chi/chi/v5"
)

func main() {
	db, err := base.CreateDB()
	// Подключаемся к БД
	if err != nil {
		log.Printf("Ошибка подключения к БД")
		fmt.Println(err)
		return

	}
	defer db.Close()

	r := chi.NewRouter()

	r.Post("/cost/", handlers.CostSummary)
	r.Put("/subscribe/{id}", handlers.HandlerUpdateSubscribe(db))

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
