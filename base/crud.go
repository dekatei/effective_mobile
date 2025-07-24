package base

import (
	"database/sql"
	"log"
	"time"
)

type Subscribe struct {
	ID        int       `json:"id"`
	UserID    string    `json:"user_id"`            //ID пользователя в формате UUID
	Service   string    `json:"service_name"`       //Название сервиса, предоставляющего подписку
	Price     int       `json:"price"`              //Стоимость месячной подписки в рублях
	StartDate time.Time `json:"start_date"`         //Дата начала подписки (месяц и год)
	EndDate   time.Time `json:"end_date,omitempty"` //Опционально дата окончания подписки
}

func InsertSubscribe(db *sql.DB, subscribe Subscribe) (int, error) {
	var id int
	err := db.QueryRow(`
		INSERT INTO subscribes (user_id, service_name, price, start_date,end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, subscribe.UserID, subscribe.Service, subscribe.Price, subscribe.StartDate, subscribe.EndDate).Scan(&id)

	return id, err
}

// SelectSubscribeByID позволяет получить подписку по номеру в таблице
func SelectSubscribeByID(db *sql.DB, id string) (Subscribe, error) {
	s := Subscribe{}

	row := db.QueryRow(`
		SELECT id, user_id, service_name, price, start_date, end_date
		FROM subscribes 
		WHERE id = $1
		`, id)

	err := row.Scan(&s.ID, &s.UserID, &s.Service, &s.Price, &s.StartDate, &s.EndDate)
	if err == sql.ErrNoRows {
		log.Printf("Подписка с ID=%d не найдена", id)
		return s, err
	} else if err != nil {
		log.Printf("Ошибка при получении подписки с ID=%d: %v", id, err)
		return s, err
	}

	return s, nil
}

// SelectUsersSubscribes позволяет получить подписки пользователя с указанием и без указания названия сервиса
func SelectUsersSubscribes(db *sql.DB, userID string, serviceName string) ([]Subscribe, error) {
	subscribes := []Subscribe{}

	query := `
		SELECT id, user_id, service_name, price, start_date, end_date
		FROM subscribes 
		WHERE user_id = $1
	`
	args := []interface{}{userID}

	if serviceName != "" {
		query += " AND service_name = $2"
		args = append(args, serviceName)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Printf("SelectSubscribeByID: query failed: %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		s := Subscribe{}

		err := rows.Scan(&s.ID, &s.UserID, &s.Service, &s.Price, &s.StartDate, &s.EndDate)
		if err != nil {
			log.Printf("SelectSubscribeByID: scan failed: %v", err)
			return nil, err
		}
		subscribes = append(subscribes, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subscribes, err
}

func DeleteSubscribe(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM subscribes WHERE id = $1", id)

	return err
}

func UpdateSubscribe(db *sql.DB, subscribe Subscribe) error {
	query := `
		UPDATE subscribes
		SET user_id = $1,
		    service_name = $2,
		    price = $3,
		    start_date = $4,
		    end_date = $5
		WHERE id = $6
	`
	_, err := db.Exec(query,
		subscribe.UserID,
		subscribe.Service,
		subscribe.Price,
		subscribe.StartDate,
		subscribe.EndDate,
		subscribe.ID,
	)

	if err != nil {
		log.Printf("UpdateSubscribe: failed to update subscription: %v", err)
		return err
	}

	return nil

}
