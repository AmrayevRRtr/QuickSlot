package handler

import (
	"QuickSlot/internal/middleware"
	"QuickSlot/internal/repository"
	"QuickSlot/internal/service"
	"encoding/json"
	"net/http"
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

	json.NewEncoder(w).Encode(appointment)
}
