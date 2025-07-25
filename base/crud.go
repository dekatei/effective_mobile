package base

import (
	"database/sql"
	"log"
	"time"
)

type Subscription struct {
	ID        int       `json:"id"`
	UserID    string    `json:"user_id"`      //ID пользователя в формате UUID
	Service   string    `json:"service_name"` //Название сервиса, предоставляющего подписку
	Price     int       `json:"price"`        //Стоимость месячной подписки в рублях
	StartDate time.Time `json:"start_date"`   //Дата начала подписки (месяц и год)
	EndDate   time.Time `json:"end_date"`     //Дата окончания подписки
}

func InsertSubscription(db *sql.DB, subscription Subscription) (int, error) {
	var id int

	err := db.QueryRow(`
		INSERT INTO subscriptions (user_id, service_name, price, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, subscription.UserID, subscription.Service, subscription.Price, subscription.StartDate, subscription.EndDate).Scan(&id)

	return id, err
}

// SelectsubscriptionByID позволяет получить подписку по номеру в таблице
func SelectSubscriptionByID(db *sql.DB, id string) (Subscription, error) {
	s := Subscription{}

	row := db.QueryRow(`
		SELECT id, user_id, service_name, price, start_date, end_date
		FROM subscriptions 
		WHERE id = $1
		`, id)

	err := row.Scan(&s.ID, &s.UserID, &s.Service, &s.Price, &s.StartDate, &s.EndDate)
	if err == sql.ErrNoRows {
		log.Printf("Подписка с ID=%d не найдена", id)
		return s, err
	} else if err != nil {
		log.Printf("Ошибка при получении подписки с ID=%s: %v", id, err)
		return s, err
	}

	return s, nil
}

// SelectUserssubscriptions позволяет получить подписки пользователя с указанием и без указания названия сервиса
func SelectUsersSubscriptions(db *sql.DB, userID string, serviceName string) ([]Subscription, error) {
	subscriptions := []Subscription{}

	query := `
		SELECT id, user_id, service_name, price, start_date, end_date
		FROM subscriptions 
		WHERE user_id = $1
	`
	args := []interface{}{userID}

	if serviceName != "" {
		query += " AND service_name = $2"
		args = append(args, serviceName)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Printf("SelectsubscriptionByID: query failed: %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		s := Subscription{}

		err := rows.Scan(&s.ID, &s.UserID, &s.Service, &s.Price, &s.StartDate, &s.EndDate)
		if err != nil {
			log.Printf("SelectsubscriptionByID: scan failed: %v", err)
			return nil, err
		}
		subscriptions = append(subscriptions, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subscriptions, err
}

func DeleteSubscription(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM subscriptions WHERE id = $1", id)

	return err
}

func UpdateSubscription(db *sql.DB, subscription Subscription) error {
	query := `
		UPDATE subscriptions
		SET user_id = $1,
		    service_name = $2,
		    price = $3,
		    start_date = $4,
		    end_date = $5
		WHERE id = $6
	`
	_, err := db.Exec(query,
		subscription.UserID,
		subscription.Service,
		subscription.Price,
		subscription.StartDate,
		subscription.EndDate,
		subscription.ID,
	)

	if err != nil {
		log.Printf("Update subscription: failed to update subscription: %v", err)
		return err
	}

	return nil

}
