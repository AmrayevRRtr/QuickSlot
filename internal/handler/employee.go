package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"QuickSlot/internal/middleware"
	"QuickSlot/internal/model"
	"QuickSlot/internal/repository"
	"QuickSlot/internal/service"
)

type EmployeeHandler struct {
	service *service.EmployeeService
}

func NewEmployeeHandler(s *service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: s}
}

type createEmployeeRequest struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	OrganizationID int64  `json:"organization_id"`
}

func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_ = ctx.Value(middleware.UserContextKey) // ensure auth

	var req createEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	if req.Email == "" && req.Phone == "" {
		http.Error(w, "email or phone is required", http.StatusBadRequest)
		return
	}

	emp := &model.Employee{
		Name:           req.Name,
		Email:          req.Email,
		Phone:          req.Phone,
		OrganizationID: req.OrganizationID,
	}

	id, err := h.service.Create(ctx, emp)
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"id": id})
}

func (h *EmployeeHandler) GetByOrganization(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgIDStr := r.URL.Query().Get("organization_id")
	if orgIDStr == "" {
		http.Error(w, "organization_id is required", http.StatusBadRequest)
		return
	}

	orgID, err := strconv.ParseInt(orgIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid organization_id", http.StatusBadRequest)
		return
	}

	employees, err := h.service.GetByOrganization(ctx, orgID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(employees)
}

func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type updateRequest struct {
		ID             int64   `json:"id"`
		Name           *string `json:"name"`
		Email          *string `json:"email"`
		Phone          *string `json:"phone"`
		OrganizationID *int64  `json:"organization_id"`
	}

	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	update := &model.EmployeeUpdate{
		Name:           req.Name,
		Email:          req.Email,
		Phone:          req.Phone,
		OrganizationID: req.OrganizationID,
	}

	if err := h.service.Update(ctx, req.ID, update); err != nil {
		if errors.Is(err, repository.ErrConflict) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{"updated": true})
}

func (h *EmployeeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		ID int64 `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(ctx, req.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
