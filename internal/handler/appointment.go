package handler

import (
	"QuickSlot/internal/middleware"
	"QuickSlot/internal/repository"
	"QuickSlot/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AppointmentHandler struct {
	service *service.AppointmentService
}

func NewAppointmentHandler(s *service.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{service: s}
}

type bookRequest struct {
	SlotID int64 `json:"slot_id"`
}

type cancelRequest struct {
	AppointmentID int64 `json:"appointment_id"`
}

func (h *AppointmentHandler) Book(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := r.Context().Value(middleware.UserContextKey).(int64)

	var req bookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	appointment, err := h.service.BookSlot(ctx, userID, req.SlotID)
	if err != nil {
		switch err {
		case repository.ErrSlotAlreadyBooked:
			http.Error(w, err.Error(), http.StatusConflict)
		case repository.ErrSlotNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	_ = json.NewEncoder(w).Encode(appointment)
}

func (h *AppointmentHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := r.Context().Value(middleware.UserContextKey).(int64)

	var req cancelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CancelBooking(ctx, userID, req.AppointmentID); err != nil {
		switch err {
		case repository.ErrAppointmentNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{
		"cancelled": true,
	})
}

func (h *AppointmentHandler) History(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := r.Context().Value(middleware.UserContextKey).(int64)

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

	items, err := h.service.GetUserHistory(ctx, userID, from, to)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(items)
}

func parseOptionalTime(v string) (*time.Time, error) {
	if v == "" {
		return nil, nil
	}

	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z0700",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, v)
		if err == nil {
			tt := t.UTC()
			return &tt, nil
		}
	}

	return nil, fmt.Errorf("invalid time format: %q", v)
}
