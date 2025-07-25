package base

import (
	"database/sql"
	"fmt"
	"time"
)

type CostFilter struct {
	UserID      string
	ServiceName string
	StartDate   time.Time
	EndDate     time.Time
}

// CountSubscriptionsCost считает суммарную стоимость подписок по фильтру
func CountSubscriptionsCost(db *sql.DB, filter CostFilter) (int, error) {
	query := `
        SELECT user_id, service_name, price, start_date, end_date
        FROM subscriptions
        WHERE (start_date <= $2) AND (end_date IS NULL OR end_date >= $1)
    `
	args := []interface{}{filter.StartDate, filter.EndDate}

	if filter.UserID != "" {
		query += " AND user_id = $3"
		args = append(args, filter.UserID)
	}
	if filter.ServiceName != "" {
		query += " AND service_name = $" + fmt.Sprint(len(args)+1)
		args = append(args, filter.ServiceName)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var total int

	for rows.Next() {
		var s Subscription
		err := rows.Scan(&s.UserID, &s.Service, &s.Price, &s.StartDate, &s.EndDate)
		if err != nil {
			return 0, err
		}

		start := maxDate(s.StartDate, filter.StartDate)
		end := minDate(s.EndDate, filter.EndDate)

		if start.After(end) {
			continue // подписка не активна в заданном периоде
		}

		months := countMonths(start, end)
		total += s.Price * months
	}

	return total, nil
}

func maxDate(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func minDate(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

// считает количество месяцев между двумя датами включительно
func countMonths(start, end time.Time) int {
	year1, month1 := start.Year(), int(start.Month())
	year2, month2 := end.Year(), int(end.Month())

	return (year2-year1)*12 + (month2 - month1) + 1
}
