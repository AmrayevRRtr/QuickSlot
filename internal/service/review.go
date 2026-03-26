package service

import (
	"QuickSlot/internal/model"
	"QuickSlot/internal/repository"
	"context"
	"errors"
)

type ReviewService struct {
	repo repository.ReviewRepository
}

func NewReviewService(r repository.ReviewRepository) *ReviewService {
	return &ReviewService{repo: r}
}

func (s *ReviewService) Create(ctx context.Context, userID, orgID int64, rating int, comment string) (int64, error) {
	if rating < 1 || rating > 5 {
		return 0, errors.New("rating must be between 1 and 5")
	}

	review := &model.Review{
		UserID:         userID,
		OrganizationID: orgID,
		Rating:         rating,
		Comment:        comment,
	}

	return s.repo.Create(ctx, review)
}

func (s *ReviewService) GetByOrganization(ctx context.Context, orgID int64) ([]model.Review, error) {
	return s.repo.GetByOrganization(ctx, orgID)
}

func (s *ReviewService) Delete(ctx context.Context, reviewID, userID int64) error {
	return s.repo.Delete(ctx, reviewID, userID)
}
