package service

import (
	"QuickSlot/internal/model"
	"QuickSlot/internal/repository"
	"context"
)

type EmployeeService struct {
	repo repository.EmployeeRepository
}

func NewEmployeeService(r repository.EmployeeRepository) *EmployeeService {
	return &EmployeeService{repo: r}
}

func (s *EmployeeService) Create(ctx context.Context, name string, orgID int64) (int64, error) {
	emp := &model.Employee{
		Name:           name,
		OrganizationID: orgID,
	}
	return s.repo.CreateEmployee(ctx, emp)
}

func (s *EmployeeService) GetByOrganization(ctx context.Context, orgID int64) ([]model.Employee, error) {
	return s.repo.GetByOrganization(ctx, orgID)
}

func (s *EmployeeService) GetByID(ctx context.Context, id int64) (*model.Employee, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *EmployeeService) Update(ctx context.Context, id int64, name string) error {
	emp := &model.Employee{ID: id, Name: name}
	return s.repo.Update(ctx, emp)
}

func (s *EmployeeService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
