package handlers

import (
	"database/sql"
	"effective_mobile/base"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type SubscriptionRequest struct {
	UserID    string `json:"user_id"`
	Service   string `json:"service_name"`
	Price     int    `json:"price"`
	StartDate string `json:"start_date"` // формат MM-YYYY
	EndDate   string `json:"end_date"`   // формат MM-YYYY
}

// HandlerAddSubscription добавляет новую подписку
// @Summary Добавить подписку
// @Description Добавляет новую подписку с указанием user_id, service_name, price и start_date
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body base.Subscription true "Данные подписки"
// @Success 200 {object} map[string]interface{} "ID новой подписки"
// @Failure 400 {object} map[string]string "Ошибка валидации или форматирования запроса"
// @Failure 405 {object} map[string]string "Метод не разрешен"
// @Failure 500 {object} map[string]string "Ошибка сервера при добавлении подписки"
// @Router /subscription [post]
func HandlerAddSubscription(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf("Получен запрос: %s %s", req.Method, req.URL.Path)

		if req.Method != http.MethodPost {
			http.Error(w, `{"error": "Метод не разрешен"}`, http.StatusMethodNotAllowed)
			return
		}

		var reqData SubscriptionRequest
		if err := json.NewDecoder(req.Body).Decode(&reqData); err != nil {
			log.Printf("Ошибка десериализации JSON в HandlerAddSubscription: %v", err)
			http.Error(w, `{"error": "Ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}

		if err := validateSubscriptionInput(reqData); err != nil {
			log.Printf("Отсутствуют обязательные поля")
			http.Error(w, `{"error": "Отсутствуют обязательные поля: user_id, service_name, price, start_date"}`, http.StatusBadRequest)
			return
		}

		parsedStartDate, err := time.Parse("01-2006", reqData.StartDate)
		if err != nil {
			log.Printf("Неправильный формат времени")
			http.Error(w, `{"error": введите время в формате 01-2006"}`, http.StatusBadRequest)
			return
		}

		var parsedEndDate time.Time
		if reqData.EndDate != "" {
			parsedEndDate, err = time.Parse("01-2006", reqData.EndDate)
			if err != nil {
				log.Printf("Неправильный формат времени")
				http.Error(w, `{"error": введите время в формате 01-2006"}`, http.StatusBadRequest)
				return
			}
		} else {
			parsedEndDate = time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC)
		}

		s := base.Subscription{
			UserID:    reqData.UserID,
			Service:   reqData.Service,
			Price:     reqData.Price,
			StartDate: parsedStartDate,
			EndDate:   parsedEndDate,
		}

		// Выполняем обновление
		id, err := base.InsertSubscription(db, s)
		if err != nil {
			log.Printf("Ошибка при добавлении подписки в InsertSubscription: %v", err)
			http.Error(w, `{"error": "Не удалось добавить подписку"}`, http.StatusInternalServerError)
			return
		}

		// Успешный ответ
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		response := map[string]interface{}{"id": id}
		json.NewEncoder(w).Encode(response)
	}
}

// HandlerAddSubscription добавляет новую подписку
// @Summary Добавить подписку
// @Description Добавляет новую подписку с указанием user_id, service_name, price и start_date
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body base.Subscription true "Данные подписки"
// @Success 200 {object} map[string]interface{} "ID новой подписки"
// @Failure 400 {object} map[string]string "Ошибка валидации или форматирования запроса"
// @Failure 405 {object} map[string]string "Метод не разрешен"
// @Failure 500 {object} map[string]string "Ошибка сервера при добавлении подписки"
// @Router /subscription [post]
func HandlerUpdateSubscription(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf("Получен запрос: %s %s", req.Method, req.URL.Path)

		if req.Method != http.MethodPut {
			http.Error(w, `{"error": "Метод не разрешен"}`, http.StatusMethodNotAllowed)
			return
		}

		// Получаем id из URL-параметров (предполагается, что используется chi или gorilla/mux)
		idStr := chi.URLParam(req, "id")
		if idStr == "" {
			http.Error(w, `{"error": "Не указан ID в URL"}`, http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			http.Error(w, `{"error": "Некорректный ID в URL"}`, http.StatusBadRequest)
			return
		}

		var reqData SubscriptionRequest
		if err := json.NewDecoder(req.Body).Decode(&reqData); err != nil {
			log.Printf("Ошибка десериализации JSON в HandlerAddSubscription: %v", err)
			http.Error(w, `{"error": "Ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}

		if err := validateSubscriptionInput(reqData); err != nil {
			log.Printf("Отсутствуют обязательные поля")
			http.Error(w, `{"error": "Отсутствуют обязательные поля: user_id, service_name, price, start_date"}`, http.StatusBadRequest)
			return
		}

		parsedStartDate, err := time.Parse("01-2006", reqData.StartDate)

		var parsedEndDate time.Time
		if reqData.EndDate != "" {
			parsedEndDate, err = time.Parse("01-2006", reqData.EndDate)
			if err != nil {
				log.Printf("Неправильный формат времени")
				http.Error(w, `{"error": введите время в формате 01-2006"}`, http.StatusBadRequest)
				return
			}
		} else {
			parsedEndDate = time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC)
		}

		if err := validateSubscriptionInput(reqData); err != nil {
			log.Printf("Отсутствуют обязательные поля")
			http.Error(w, `{"error": "Отсутствуют обязательные поля: user_id, service_name, price, start_date"}`, http.StatusBadRequest)
			return
		}

		s := base.Subscription{
			UserID:    reqData.UserID,
			Service:   reqData.Service,
			Price:     reqData.Price,
			StartDate: parsedStartDate,
			EndDate:   parsedEndDate,
		}

		// Принудительно устанавливаем ID из URL (игнорируем ID из тела)
		s.ID = id
		// Проверки
		if s.ID == 0 {
			log.Printf("Не указан ID подписки")
			http.Error(w, `{"error": "Не указан ID подписки"}`, http.StatusBadRequest)
			return
		}

		// Выполняем обновление
		err = base.UpdateSubscription(db, s)
		if err != nil {
			log.Printf("Ошибка при обновлении подписки: %v", err)
			http.Error(w, `{"error": "Не удалось обновить подписку"}`, http.StatusInternalServerError)
			return
		}

		// Успешный ответ
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}

// HandlerDeleteSubscription удаляет подписку по ID
// @Summary Удалить подписку по ID
// @Description Удаляет подписку по указанному ID
// @Tags subscriptions
// @Produce json
// @Param id path int true "ID подписки"
// @Success 200 {object} map[string]interface{} "ID удаленной подписки"
// @Failure 405 {object} map[string]string "Метод не разрешен"
// @Failure 500 {object} map[string]string "Ошибка сервера при удалении подписки"
// @Router /subscription/{id} [delete]
func HandlerDeleteSubscription(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf("Получен запрос: %s %s", req.Method, req.URL.Path)

		if req.Method != http.MethodDelete {
			http.Error(w, `{"error": "Метод не разрешен"}`, http.StatusMethodNotAllowed)
			return
		}

		idStr := chi.URLParam(req, "id")

		// Выполняем удаление
		err := base.DeleteSubscription(db, idStr)
		if err != nil {
			log.Printf("Ошибка при удалении подписки в DeleteSubscription: %v", err)
			http.Error(w, `{"error": "Не удалось удалить подписку"}`, http.StatusInternalServerError)
			return
		}

		// Успешный ответ
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		response := map[string]interface{}{"id": idStr}
		json.NewEncoder(w).Encode(response)
	}
}

func validateSubscriptionInput(s SubscriptionRequest) error {
	if s.UserID == "" || s.Service == "" || s.Price <= 0 || s.StartDate == "" {
		return errors.New("user_id, service_name, price и start_date обязательны")
	}
	return nil
}

// HandlerGetSubscriptionsByUserID возвращает список подписок пользователя (с фильтром по названию сервиса)
// @Summary Получить подписки пользователя
// @Description Получить все подписки по user_id, опционально фильтруя по service_name
// @Tags subscriptions
// @Produce json
// @Param user_id path string true "ID пользователя"
// @Param service_name query string false "Название сервиса (опционально)"
// @Success 200 {array} base.Subscription
// @Failure 400 {object} map[string]string "Ошибка валидации запроса"
// @Failure 404 {object} map[string]string "Подписки не найдены"
// @Failure 500 {object} map[string]string "Ошибка сервера при получении подписок"
// @Router /subscriptions/{user_id} [get]
func HandlerGetSubscriptionsByUserID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Получен запрос: %s %s", r.Method, r.URL.Path)

		userID := chi.URLParam(r, "user_id")
		if userID == "" {
			http.Error(w, `{"error": "Не указан user_id в URL"}`, http.StatusBadRequest)
			return
		}
		serviceName := r.URL.Query().Get("service_name")
		subs, err := base.SelectUsersSubscriptions(db, userID, serviceName)
		if err != nil {
			log.Printf("Ошибка при получении подписок: %v", err)
			http.Error(w, `{"error": "Ошибка при получении подписок"}`, http.StatusInternalServerError)
			return
		}

		if len(subs) == 0 {
			http.Error(w, `{"error": "Подписки не найдены"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(subs)
	}
}

// HandlerGetSubscriptionByID возвращает подписку по ID
// @Summary Получить подписку по ID
// @Description Возвращает подписку по уникальному ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {object} base.Subscription
// @Failure 400 {object} map[string]string "Ошибка валидации запроса"
// @Failure 500 {object} map[string]string "Ошибка сервера при получении подписки"
// @Router /subscription/{id} [get]
func HandlerGetSubscriptionByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Получен запрос: %s %s", r.Method, r.URL.Path)

		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, `{"error": "Не указан id в URL"}`, http.StatusBadRequest)
			return
		}

		sub, err := base.SelectSubscriptionByID(db, id)
		if err != nil {
			log.Printf("Ошибка при получении подписки: %v", err)
			http.Error(w, `{"error": "Ошибка при получении подписки"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(sub)
	}
}
