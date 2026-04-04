package handler

import (
	"QuickSlot/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type SlotHandler struct {
	service *service.SlotService
}

func NewSlotHandler(s *service.SlotService) *SlotHandler {
	return &SlotHandler{service: s}
}

type generateRequest struct {
	EmployeeID int64  `json:"employee_id"`
	Date       string `json:"date"`
	StartHour  int    `json:"start_hour"`
	EndHour    int    `json:"end_hour"`
	Duration   int    `json:"duration"`
}

func (h *SlotHandler) Generate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req generateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "invalid date format", http.StatusBadRequest)
		return
	}

	err = h.service.GenerateSlots(
		ctx,
		req.EmployeeID,
		date,
		req.StartHour,
		req.EndHour,
		req.Duration,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte("slots generated"))
}

func (h *SlotHandler) GetAvailableByEmployee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	employeeIDStr := r.URL.Query().Get("employee_id")
	if employeeIDStr == "" {
		http.Error(w, "employee_id is required", http.StatusBadRequest)
		return
	}

	employeeID, err := strconv.ParseInt(employeeIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid employee_id", http.StatusBadRequest)
		return
	}

	from, err := parseOptionalTime(r.URL.Query().Get("from"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	to, err := parseOptionalTime(r.URL.Query().Get("to"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	slots, err := h.service.GetAvailableByEmployee(ctx, employeeID, from, to)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(slots)
}
