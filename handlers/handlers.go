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

type SubscribeRequest struct {
	UserID    string `json:"user_id"`
	Service   string `json:"service_name"`
	Price     int    `json:"price"`
	StartDate string `json:"start_date"` // формат MM-YYYY
	EndDate   string `json:"end_date"`   // формат MM-YYYY
}

// HandlerAddSubscribe добавляет новую подписку
// @Summary Добавить подписку
// @Description Добавляет новую подписку с указанием user_id, service_name, price и start_date
// @Tags subscribes
// @Accept json
// @Produce json
// @Param subscribe body base.Subscribe true "Данные подписки"
// @Success 200 {object} map[string]interface{} "ID новой подписки"
// @Failure 400 {object} map[string]string "Ошибка валидации или форматирования запроса"
// @Failure 405 {object} map[string]string "Метод не разрешен"
// @Failure 500 {object} map[string]string "Ошибка сервера при добавлении подписки"
// @Router /subscribe [post]
func HandlerAddSubscribe(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf("Получен запрос: %s %s", req.Method, req.URL.Path)

		if req.Method != http.MethodPost {
			http.Error(w, `{"error": "Метод не разрешен"}`, http.StatusMethodNotAllowed)
			return
		}

		var reqData SubscribeRequest
		if err := json.NewDecoder(req.Body).Decode(&reqData); err != nil {
			log.Printf("Ошибка десериализации JSON в HandlerAddSubscribe: %v", err)
			http.Error(w, `{"error": "Ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}

		if err := validateSubscribeInput(reqData); err != nil {
			log.Printf("Отсутствуют обязательные поля")
			http.Error(w, `{"error": "Отсутствуют обязательные поля: user_id, service_name, price, start_date"}`, http.StatusBadRequest)
			return
		}

		parsedStartDate, err := time.Parse("01-2006", reqData.StartDate)

		var parsedEndDate time.Time
		if reqData.EndDate != "" {
			parsedEndDate, err = time.Parse("01-2006", reqData.EndDate)
		}
		s := base.Subscribe{
			UserID:    reqData.UserID,
			Service:   reqData.Service,
			Price:     reqData.Price,
			StartDate: parsedStartDate,
			EndDate:   parsedEndDate,
		}

		// Выполняем обновление
		id, err := base.InsertSubscribe(db, s)
		if err != nil {
			log.Printf("Ошибка при добавлении подписки в InsertSubscribe: %v", err)
			http.Error(w, `{"error": "Не удалось добавить подписку"}`, http.StatusInternalServerError)
			return
		}

		// Успешный ответ
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		response := map[string]interface{}{"id": id}
		json.NewEncoder(w).Encode(response)
	}
}

// HandlerAddSubscribe добавляет новую подписку
// @Summary Добавить подписку
// @Description Добавляет новую подписку с указанием user_id, service_name, price и start_date
// @Tags subscribes
// @Accept json
// @Produce json
// @Param subscribe body base.Subscribe true "Данные подписки"
// @Success 200 {object} map[string]interface{} "ID новой подписки"
// @Failure 400 {object} map[string]string "Ошибка валидации или форматирования запроса"
// @Failure 405 {object} map[string]string "Метод не разрешен"
// @Failure 500 {object} map[string]string "Ошибка сервера при добавлении подписки"
// @Router /subscribe [post]
func HandlerUpdateSubscribe(db *sql.DB) http.HandlerFunc {
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

		var reqData SubscribeRequest
		if err := json.NewDecoder(req.Body).Decode(&reqData); err != nil {
			log.Printf("Ошибка десериализации JSON в HandlerAddSubscribe: %v", err)
			http.Error(w, `{"error": "Ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}

		if err := validateSubscribeInput(reqData); err != nil {
			log.Printf("Отсутствуют обязательные поля")
			http.Error(w, `{"error": "Отсутствуют обязательные поля: user_id, service_name, price, start_date"}`, http.StatusBadRequest)
			return
		}

		parsedStartDate, err := time.Parse("01-2006", reqData.StartDate)

		var parsedEndDate time.Time
		if reqData.EndDate != "" {
			parsedEndDate, err = time.Parse("01-2006", reqData.EndDate)
		}

		if err := validateSubscribeInput(reqData); err != nil {
			log.Printf("Отсутствуют обязательные поля")
			http.Error(w, `{"error": "Отсутствуют обязательные поля: user_id, service_name, price, start_date"}`, http.StatusBadRequest)
			return
		}

		s := base.Subscribe{
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
		err = base.UpdateSubscribe(db, s)
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

// HandlerDeleteSubscribe удаляет подписку по ID
// @Summary Удалить подписку по ID
// @Description Удаляет подписку по указанному ID
// @Tags subscribes
// @Produce json
// @Param id path int true "ID подписки"
// @Success 200 {object} map[string]interface{} "ID удаленной подписки"
// @Failure 405 {object} map[string]string "Метод не разрешен"
// @Failure 500 {object} map[string]string "Ошибка сервера при удалении подписки"
// @Router /subscribe/{id} [delete]
func HandlerDeleteSubscribe(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf("Получен запрос: %s %s", req.Method, req.URL.Path)

		if req.Method != http.MethodDelete {
			http.Error(w, `{"error": "Метод не разрешен"}`, http.StatusMethodNotAllowed)
			return
		}

		idStr := chi.URLParam(req, "id")

		// Выполняем удаление
		err := base.DeleteSubscribe(db, idStr)
		if err != nil {
			log.Printf("Ошибка при удалении подписки в DeleteSubscribe: %v", err)
			http.Error(w, `{"error": "Не удалось удалить подписку"}`, http.StatusInternalServerError)
			return
		}

		// Успешный ответ
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		response := map[string]interface{}{"id": idStr}
		json.NewEncoder(w).Encode(response)
	}
}

func validateSubscribeInput(s SubscribeRequest) error {
	if s.UserID == "" || s.Service == "" || s.Price <= 0 || s.StartDate == "" {
		return errors.New("user_id, service_name, price и start_date обязательны")
	}
	return nil
}

// HandlerGetSubscribesByUserID возвращает список подписок пользователя (с фильтром по названию сервиса)
// @Summary Получить подписки пользователя
// @Description Получить все подписки по user_id, опционально фильтруя по service_name
// @Tags subscribes
// @Produce json
// @Param user_id path string true "ID пользователя"
// @Param service_name query string false "Название сервиса (опционально)"
// @Success 200 {array} base.Subscribe
// @Failure 400 {object} map[string]string "Ошибка валидации запроса"
// @Failure 404 {object} map[string]string "Подписки не найдены"
// @Failure 500 {object} map[string]string "Ошибка сервера при получении подписок"
// @Router /subscribes/{user_id} [get]
func HandlerGetSubscribesByUserID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Получен запрос: %s %s", r.Method, r.URL.Path)

		userID := chi.URLParam(r, "user_id")
		if userID == "" {
			http.Error(w, `{"error": "Не указан user_id в URL"}`, http.StatusBadRequest)
			return
		}
		serviceName := r.URL.Query().Get("service_name")
		subs, err := base.SelectUsersSubscribes(db, userID, serviceName)
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

// HandlerGetSubscribeByID возвращает подписку по ID
// @Summary Получить подписку по ID
// @Description Возвращает подписку по уникальному ID
// @Tags subscribes
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {object} base.Subscribe
// @Failure 400 {object} map[string]string "Ошибка валидации запроса"
// @Failure 500 {object} map[string]string "Ошибка сервера при получении подписки"
// @Router /subscribe/{id} [get]
func HandlerGetSubscribeByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Получен запрос: %s %s", r.Method, r.URL.Path)

		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, `{"error": "Не указан id в URL"}`, http.StatusBadRequest)
			return
		}

		sub, err := base.SelectSubscribeByID(db, id)
		if err != nil {
			log.Printf("Ошибка при получении подписки: %v", err)
			http.Error(w, `{"error": "Ошибка при получении подписки"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(sub)
	}
}
