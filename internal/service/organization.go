package service

import (
	"context"

	"QuickSlot/internal/model"
	"QuickSlot/internal/repository"
)

type OrganizationService struct {
	repo repository.OrganizationRepository
}

func NewOrganizationService(r repository.OrganizationRepository) *OrganizationService {
	return &OrganizationService{repo: r}
}

func (s *OrganizationService) Create(ctx context.Context, name string, ownerID int64) (int64, error) {
	exists, err := s.repo.ExistsByName(ctx, name)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, repository.ErrConflict
	}

	org := &model.Organization{
		Name:    name,
		OwnerID: ownerID,
	}
	return s.repo.CreateOrg(ctx, org)
}

func (s *OrganizationService) GetAll(ctx context.Context) ([]model.Organization, error) {
	return s.repo.GetAll(ctx)
}

func (s *OrganizationService) GetByID(ctx context.Context, id int64) (*model.Organization, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *OrganizationService) Update(ctx context.Context, id int64, update *model.OrganizationUpdate) error {
	return s.repo.Update(ctx, id, update)
}

func (s *OrganizationService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
