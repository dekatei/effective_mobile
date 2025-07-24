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
		INSERT INTO subscribes (user_id, service_name, price, start_date, expiration_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, subscribe.UserID, subscribe.Service, subscribe.Price, subscribe.StartDate, subscribe.EndDate).Scan(&id)

	return id, err
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
		log.Printf("SelectUsersSubscribes: query failed: %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		s := Subscribe{}

		err := rows.Scan(&s.ID, &s.UserID, &s.Service, &s.Price, &s.StartDate, &s.EndDate)
		if err != nil {
			log.Printf("SelectUsersSubscribes: scan failed: %v", err)
			return nil, err
		}
		subscribes = append(subscribes, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subscribes, err
}

func DeleteSubscribe(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM subscribes WHERE id = ?", id)

	return err
}
