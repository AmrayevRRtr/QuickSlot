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

type OrganizationHandler struct {
	service *service.OrganizationService
}

func NewOrganizationHandler(s *service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{service: s}
}

type createOrgRequest struct {
	Name string `json:"name"`
}

func (h *OrganizationHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(middleware.UserContextKey).(int64)

	var req createOrgRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(ctx, req.Name, userID)
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"id": id})
}

func (h *OrganizationHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgs, err := h.service.GetAll(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(orgs)
}

func (h *OrganizationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	org, err := h.service.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrOrgNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(org)
}

func (h *OrganizationHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type updateRequest struct {
		ID   int64   `json:"id"`
		Name *string `json:"name"`
	}

	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	update := &model.OrganizationUpdate{
		Name: req.Name,
	}

	if err := h.service.Update(ctx, req.ID, update); err != nil {
		if errors.Is(err, repository.ErrConflict) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{"updated": true})
}

func (h *OrganizationHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
