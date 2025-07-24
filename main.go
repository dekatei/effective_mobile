package main

import (
	"net/http"

	"effective_mobile/handlers"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Post("/cost/", handlers.CostSummary)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
