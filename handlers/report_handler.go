package handlers

import (
	"encoding/json"
	"net/http"
	"product-api/services"
	"strings"
	"time"
)

const reportDateLayout = "2006-01-02"

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) HandleTodayReport(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetToday(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ReportHandler) HandleReport(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetByRange(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ReportHandler) GetToday(w http.ResponseWriter, r *http.Request) {
	startInclusive, endExclusive := dayRange(time.Now())

	report, err := h.service.GetReport(startInclusive, endExclusive)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

func (h *ReportHandler) GetByRange(w http.ResponseWriter, r *http.Request) {
	startDateStr := strings.TrimSpace(r.URL.Query().Get("start_date"))
	endDateStr := strings.TrimSpace(r.URL.Query().Get("end_date"))
	if startDateStr == "" || endDateStr == "" {
		http.Error(w, "start_date and end_date are required", http.StatusBadRequest)
		return
	}

	loc := time.Local
	startDate, err := parseISODate(startDateStr, loc)
	if err != nil {
		http.Error(w, "Invalid start_date", http.StatusBadRequest)
		return
	}

	endDate, err := parseISODate(endDateStr, loc)
	if err != nil {
		http.Error(w, "Invalid end_date", http.StatusBadRequest)
		return
	}

	if endDate.Before(startDate) {
		http.Error(w, "end_date must be >= start_date", http.StatusBadRequest)
		return
	}

	startInclusive := startDate
	endExclusive := endDate.AddDate(0, 0, 1)

	report, err := h.service.GetReport(startInclusive, endExclusive)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

func parseISODate(value string, loc *time.Location) (time.Time, error) {
	parsed, err := time.ParseInLocation(reportDateLayout, strings.TrimSpace(value), loc)
	if err != nil {
		return time.Time{}, err
	}
	y, m, d := parsed.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, loc), nil
}

func dayRange(now time.Time) (time.Time, time.Time) {
	y, m, d := now.Date()
	loc := now.Location()
	start := time.Date(y, m, d, 0, 0, 0, 0, loc)
	return start, start.AddDate(0, 0, 1)
}
