package service

import (
	"context"
	"errors"

	"QuickSlot/internal/model"
	"QuickSlot/internal/repository"
)

type EmployeeService struct {
	empRepo repository.EmployeeRepository
	orgRepo repository.OrganizationRepository
}

func NewEmployeeService(e repository.EmployeeRepository, o repository.OrganizationRepository) *EmployeeService {
	return &EmployeeService{empRepo: e, orgRepo: o}
}

func (s *EmployeeService) Create(ctx context.Context, emp *model.Employee) (int64, error) {
	exists, err := s.empRepo.ExistsByEmailOrPhone(ctx, emp.Email, emp.Phone)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, repository.ErrConflict
	}

	_, err = s.orgRepo.GetByID(ctx, emp.OrganizationID)
	if err != nil {
		if errors.Is(err, repository.ErrOrgNotFound) {
			return 0, errors.New("organization does not exist or is deleted")
		}
		return 0, err
	}

	return s.empRepo.CreateEmployee(ctx, emp)
}

func (s *EmployeeService) GetByOrganization(ctx context.Context, orgID int64) ([]model.Employee, error) {
	return s.empRepo.GetByOrganization(ctx, orgID)
}

func (s *EmployeeService) GetByID(ctx context.Context, id int64) (*model.Employee, error) {
	return s.empRepo.GetByID(ctx, id)
}

func (s *EmployeeService) Update(ctx context.Context, id int64, update *model.EmployeeUpdate) error {
	if update.OrganizationID != nil {
		_, err := s.orgRepo.GetByID(ctx, *update.OrganizationID)
		if err != nil {
			if errors.Is(err, repository.ErrOrgNotFound) {
				return errors.New("organization does not exist or is deleted")
			}
			return err
		}
	}
	return s.empRepo.Update(ctx, id, update)
}

func (s *EmployeeService) Delete(ctx context.Context, id int64) error {
	return s.empRepo.Delete(ctx, id)
}
