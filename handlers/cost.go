package handlers

import (
	"database/sql"
	"effective_mobile/base"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// CostRequest представляет запрос на подсчет стоимости подписок
type CostRequest struct {
	// ID пользователя
	UserID string `json:"user_id,omitempty"`
	// Название сервиса (необязательно)
	ServiceName string `json:"service_name,omitempty"`
	// Дата начала в формате YYYY-MM
	StartDate string `json:"start_date" example:"01-2024"`
	// Дата окончания в формате YYYY-MM
	EndDate string `json:"end_date" example:"07-2024"`
}

// CostResponse представляет ответ с суммарной стоимостью
type CostResponse struct {
	// Общая сумма стоимости подписок
	Total int `json:"total" example:"1500"`
}

// CostSummary возвращает суммарную стоимость подписок за период для пользователя
// @Summary Суммарная стоимость подписок
// @Description Возвращает сумму подписок пользователя за указанный период с возможностью фильтрации по названию сервиса
// @Tags cost
// @Produce json
// @Param user_id path string true "ID пользователя"
// @Param start query string true "Дата начала (формат YYYY-MM)"
// @Param end query string true "Дата окончания (формат YYYY-MM)"
// @Param service_name query string false "Название сервиса (опционально)"
// @Success 200 {object} CostResponse "Общий итог по подпискам"
// @Failure 400 {object} map[string]string "Ошибки валидации параметров запроса"
// @Failure 500 {object} map[string]string "Ошибка сервера при подсчёте стоимости"
// @Router /cost/{user_id} [get]
func CostSummary(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "user_id")
		if userID == "" {
			http.Error(w, `{"error": "Не указан user_id в URL"}`, http.StatusBadRequest)
			return
		}

		// Читаем query параметры
		service := r.URL.Query().Get("service_name")
		startStr := r.URL.Query().Get("start_date")
		endStr := r.URL.Query().Get("end_date")

		if startStr == "" || endStr == "" {
			http.Error(w, `{"error": "start_date и end_date обязательны. Формат: MM-YYYY"}`, http.StatusBadRequest)
			return
		}

		startDate, err := time.Parse("01-2006", startStr)
		if err != nil {
			http.Error(w, `{"error": "Неверный формат start_date. Используйте MM-YYYY"}`, http.StatusBadRequest)
			return
		}
		endDate, err := time.Parse("01-2006", endStr)
		if err != nil {
			http.Error(w, `{"error": "Неверный формат end_date. Используйте MM-YYYY"}`, http.StatusBadRequest)
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
