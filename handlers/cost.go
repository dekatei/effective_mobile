package handlers

import (
	"database/sql"
	"effective_mobile/base"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type CostRequest struct {
	UserID      string `json:"user_id,omitempty"`
	ServiceName string `json:"service_name,omitempty"`
	StartDate   string `json:"start_date"` // формат YYYY-MM
	EndDate     string `json:"end_date"`   // формат YYYY-MM
}

type CostResponse struct {
	Total int `json:"total"`
}

func CostSummary(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "user_id")
		if userID == "" {
			http.Error(w, `{"error": "Не указан user_id в URL"}`, http.StatusBadRequest)
			return
		}

		// Читаем query параметры
		service := r.URL.Query().Get("service_name")
		startStr := r.URL.Query().Get("start")
		endStr := r.URL.Query().Get("end")

		if startStr == "" || endStr == "" {
			http.Error(w, `{"error": "start и end обязательны. Формат: YYYY-MM"}`, http.StatusBadRequest)
			return
		}

		startDate, err := time.Parse("2006-01", startStr)
		if err != nil {
			http.Error(w, `{"error": "Неверный формат start. Используйте YYYY-MM"}`, http.StatusBadRequest)
			return
		}
		endDate, err := time.Parse("2006-01", endStr)
		if err != nil {
			http.Error(w, `{"error": "Неверный формат end. Используйте YYYY-MM"}`, http.StatusBadRequest)
			return
		}

		// Собираем фильтр
		filter := base.CostFilter{
			UserID:      userID,
			ServiceName: service,
			StartDate:   startDate,
			EndDate:     endDate,
		}

		// Запрашиваем сумму
		total, err := base.CountSubscriptionsCost(db, filter)
		if err != nil {
			http.Error(w, `{"error": "Ошибка при обращении к базе данных"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(CostResponse{Total: total})
	}
}
