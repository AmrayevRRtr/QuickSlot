package service

import (
	"QuickSlot/internal/model"
	"context"
	"testing"
)

type mockReviewRepo struct {
	reviews []model.Review
	lastID  int64
}

func (m *mockReviewRepo) Create(ctx context.Context, review *model.Review) (int64, error) {
	m.lastID++
	review.ID = m.lastID
	m.reviews = append(m.reviews, *review)
	return m.lastID, nil
}

func (m *mockReviewRepo) GetByOrganization(ctx context.Context, orgID int64) ([]model.Review, error) {
	var result []model.Review
	for _, r := range m.reviews {
		if r.OrganizationID == orgID {
			result = append(result, r)
		}
	}
	return result, nil
}

func (m *mockReviewRepo) GetByID(ctx context.Context, id int64) (*model.Review, error) {
	for _, r := range m.reviews {
		if r.ID == id {
			return &r, nil
		}
	}
	return nil, nil
}

func (m *mockReviewRepo) Delete(ctx context.Context, id, userID int64) error {
	for i, r := range m.reviews {
		if r.ID == id && r.UserID == userID {
			m.reviews = append(m.reviews[:i], m.reviews[i+1:]...)
			return nil
		}
	}
	return nil
}

func TestCreateReview(t *testing.T) {
	repo := &mockReviewRepo{}
	svc := NewReviewService(repo)

	id, err := svc.Create(context.Background(), 1, 1, 5, "great service")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != 1 {
		t.Fatalf("expected id=1, got %d", id)
	}
}

func TestCreateReviewInvalidRating(t *testing.T) {
	repo := &mockReviewRepo{}
	svc := NewReviewService(repo)

	_, err := svc.Create(context.Background(), 1, 1, 0, "bad")
	if err == nil {
		t.Fatal("expected error for rating 0")
	}

	_, err = svc.Create(context.Background(), 1, 1, 6, "too high")
	if err == nil {
		t.Fatal("expected error for rating 6")
	}
}

func TestGetByOrganization(t *testing.T) {
	repo := &mockReviewRepo{}
	svc := NewReviewService(repo)

	_, _ = svc.Create(context.Background(), 1, 10, 4, "good")
	_, _ = svc.Create(context.Background(), 2, 10, 3, "ok")
	_, _ = svc.Create(context.Background(), 3, 20, 5, "other org")

	reviews, err := svc.GetByOrganization(context.Background(), 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(reviews) != 2 {
		t.Fatalf("expected 2 reviews for org 10, got %d", len(reviews))
	}
}

func TestDeleteReview(t *testing.T) {
	repo := &mockReviewRepo{}
	svc := NewReviewService(repo)

	id, _ := svc.Create(context.Background(), 1, 1, 5, "will delete")

	err := svc.Delete(context.Background(), id, 1)
	if err != nil {
		t.Fatal(err)
	}

	reviews, _ := svc.GetByOrganization(context.Background(), 1)
	if len(reviews) != 0 {
		t.Fatal("review should have been deleted")
	}
}
