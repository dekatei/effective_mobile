package handlers

import (
	"database/sql"
	"effective_mobile/base"
	"encoding/json"
	"net/http"
	"time"
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

var DB *sql.DB

func CostSummary(w http.ResponseWriter, r *http.Request) {
	var req CostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
	}
	if req.StartDate == "" || req.EndDate == "" {
		http.Error(w, "start_date and end_date are required", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01", req.StartDate)
	if err != nil {
		http.Error(w, "invalid start_date format, use YYYY-MM", http.StatusBadRequest)
		return
	}
	endDate, err := time.Parse("2006-01", req.EndDate)
	if err != nil {
		http.Error(w, "invalid end_date format, use YYYY-MM", http.StatusBadRequest)
		return
	}

	filter := base.CostFilter{
		UserID:      req.UserID,
		ServiceName: req.ServiceName,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	total, err := base.CountSubscriptionsCost(DB, filter)
	if err != nil {
		http.Error(w, "database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(CostResponse{Total: total})

}
