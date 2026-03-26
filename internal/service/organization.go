package service

import (
	"QuickSlot/internal/model"
	"QuickSlot/internal/repository"
	"context"
)

type OrganizationService struct {
	repo repository.OrganizationRepository
}

func NewOrganizationService(r repository.OrganizationRepository) *OrganizationService {
	return &OrganizationService{repo: r}
}

func (s *OrganizationService) Create(ctx context.Context, name string, ownerID int64) (int64, error) {
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

func (s *OrganizationService) Update(ctx context.Context, id int64, name string) error {
	org := &model.Organization{ID: id, Name: name}
	return s.repo.Update(ctx, org)
}

func (s *OrganizationService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
