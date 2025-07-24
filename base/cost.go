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
	query := `SELECT SUM(price) FROM subscribes WHERE start_date BETWEEN $1 AND $2`
	args := []interface{}{filter.StartDate, filter.EndDate}
	argIdx := 3

	if filter.UserID != "" {
		query += fmt.Sprintf(" AND user_id = $%d", argIdx)
		args = append(args, filter.UserID)
		argIdx++
	}
	if filter.ServiceName != "" {
		query += fmt.Sprintf(" AND service_name = $%d", argIdx)
		args = append(args, filter.ServiceName)
	}

	var total sql.NullInt64
	err := db.QueryRow(query, args...).Scan(&total)
	if err != nil {
		return 0, err
	}

	if !total.Valid {
		return 0, nil
	}
	return int(total.Int64), nil
}
