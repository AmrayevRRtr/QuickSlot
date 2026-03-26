package repository

import (
	"QuickSlot/internal/model"
	"QuickSlot/pkg/database/mysql"
	"context"
)

type ReviewRepository interface {
	Create(ctx context.Context, review *model.Review) (int64, error)
	GetByOrganization(ctx context.Context, orgID int64) ([]model.Review, error)
	GetByID(ctx context.Context, id int64) (*model.Review, error)
	Delete(ctx context.Context, id, userID int64) error
}

type reviewRepository struct {
	db *mysql.Dialect
}

func NewReviewRepository(db *mysql.Dialect) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(ctx context.Context, review *model.Review) (int64, error) {
	result, err := r.db.DB.ExecContext(ctx,
		"INSERT INTO reviews (user_id, organization_id, rating, comment) VALUES (?, ?, ?, ?)",
		review.UserID, review.OrganizationID, review.Rating, review.Comment,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (r *reviewRepository) GetByOrganization(ctx context.Context, orgID int64) ([]model.Review, error) {
	reviews := make([]model.Review, 0)
	err := r.db.DB.SelectContext(ctx, &reviews,
		"SELECT * FROM reviews WHERE organization_id = ? ORDER BY created_at DESC", orgID)
	return reviews, err
}

func (r *reviewRepository) GetByID(ctx context.Context, id int64) (*model.Review, error) {
	var review model.Review
	err := r.db.DB.GetContext(ctx, &review, "SELECT * FROM reviews WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) Delete(ctx context.Context, id, userID int64) error {
	res, err := r.db.DB.ExecContext(ctx, "DELETE FROM reviews WHERE id = ? AND user_id = ?", id, userID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrReviewNotFound
	}
	return nil
}
