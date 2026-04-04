package handler

import (
	"QuickSlot/internal/middleware"
	"QuickSlot/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
)

type ReviewHandler struct {
	service *service.ReviewService
}

func NewReviewHandler(s *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{service: s}
}

type createReviewRequest struct {
	OrganizationID int64  `json:"organization_id"`
	Rating         int    `json:"rating"`
	Comment        string `json:"comment"`
}

func (h *ReviewHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(middleware.UserContextKey).(int64)

	var req createReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(ctx, userID, req.OrganizationID, req.Rating, req.Comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"review_id": id,
	})
}

func (h *ReviewHandler) GetByOrganization(w http.ResponseWriter, r *http.Request) {
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

	reviews, err := h.service.GetByOrganization(ctx, orgID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(reviews)
}

type deleteReviewRequest struct {
	ReviewID int64 `json:"review_id"`
}

func (h *ReviewHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(middleware.UserContextKey).(int64)

	var req deleteReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(ctx, req.ReviewID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{
		"deleted": true,
	})
}
