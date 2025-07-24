package handlers

import (
	"database/sql"
	"effective_mobile/base"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// HandlerAddSubscribe обрабатывает POST-запросы для добавления подписки
func HandlerAddSubscribe(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf("Получен запрос: %s %s", req.Method, req.URL.Path)

		if req.Method != http.MethodPost {
			http.Error(w, `{"error": "Метод не разрешен"}`, http.StatusMethodNotAllowed)
			return
		}

		var s base.Subscribe
		if err := json.NewDecoder(req.Body).Decode(&s); err != nil {
			log.Printf("Ошибка десериализации JSON в HandlerAddSubscribe: %v", err)
			http.Error(w, `{"error": "Ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}

		if err := validateSubscribeInput(s); err != nil {
			log.Printf("Отсутствуют обязательные поля")
			http.Error(w, `{"error": "Отсутствуют обязательные поля: user_id, service_name, price, start_date"}`, http.StatusBadRequest)
			return
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

// HandlerUpdateSubscribe обрабатывает PUT-запросы для редактирования подписки
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

		var s base.Subscribe
		if err := json.NewDecoder(req.Body).Decode(&s); err != nil {
			log.Printf("Ошибка десериализации JSON в HandlerUpdateSubscribe: %v", err)
			http.Error(w, `{"error": "Ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}
		// Принудительно устанавливаем ID из URL (игнорируем ID из тела)
		s.ID = id
		// Проверки
		if s.ID == 0 {
			log.Printf("Не указан ID подписки")
			http.Error(w, `{"error": "Не указан ID подписки"}`, http.StatusBadRequest)
			return
		}
		if err := validateSubscribeInput(s); err != nil {
			log.Printf("Отсутствуют обязательные поля")
			http.Error(w, `{"error": "Отсутствуют обязательные поля: user_id, service_name, price, start_date"}`, http.StatusBadRequest)
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
			log.Printf("Ошибка при удалении подписки в InsertSubscribe: %v", err)
			http.Error(w, `{"error": "Не удалось удалить подписку"}`, http.StatusInternalServerError)
			return
		}

		// Успешный ответ
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		response := map[string]interface{}{"id": idStr}
		json.NewEncoder(w).Encode(response)
	}
}

func validateSubscribeInput(s base.Subscribe) error {
	if s.UserID == "" || s.Service == "" || s.Price <= 0 || s.StartDate.IsZero() {
		return errors.New("user_id, service_name, price и start_date обязательны")
	}
	return nil
}

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
