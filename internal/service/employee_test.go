package service

import (
	"QuickSlot/internal/model"
	"QuickSlot/internal/repository"
	"context"
	"errors"
	"testing"
)

type mockEmpRepo struct {
	emps   map[int64]*model.Employee
	lastID int64
}

func newMockEmpRepo() *mockEmpRepo {
	return &mockEmpRepo{emps: make(map[int64]*model.Employee)}
}

func (m *mockEmpRepo) CreateEmployee(ctx context.Context, emp *model.Employee) (int64, error) {
	m.lastID++
	emp.ID = m.lastID
	m.emps[m.lastID] = emp
	return m.lastID, nil
}

func (m *mockEmpRepo) GetByOrganization(ctx context.Context, orgID int64) ([]model.Employee, error) {
	var result []model.Employee
	for _, e := range m.emps {
		if e.OrganizationID == orgID {
			result = append(result, *e)
		}
	}
	return result, nil
}

func (m *mockEmpRepo) GetByID(ctx context.Context, id int64) (*model.Employee, error) {
	emp, exists := m.emps[id]
	if !exists {
		return nil, repository.ErrEmployeeNotFound
	}
	return emp, nil
}

func (m *mockEmpRepo) ExistsByEmailOrPhone(ctx context.Context, email, phone string) (bool, error) {
	for _, e := range m.emps {
		if (email != "" && e.Email == email) || (phone != "" && e.Phone == phone) {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockEmpRepo) Update(ctx context.Context, id int64, update *model.EmployeeUpdate) error {
	emp, exists := m.emps[id]
	if !exists {
		return repository.ErrEmployeeNotFound
	}
	if update.Name != nil {
		emp.Name = *update.Name
	}
	if update.Email != nil {
		emp.Email = *update.Email
	}
	if update.Phone != nil {
		emp.Phone = *update.Phone
	}
	if update.OrganizationID != nil {
		emp.OrganizationID = *update.OrganizationID
	}
	return nil
}

func (m *mockEmpRepo) Delete(ctx context.Context, id int64) error {
	delete(m.emps, id)
	return nil
}

func TestEmployeeCreateSuccess(t *testing.T) {
	empRepo := newMockEmpRepo()
	orgRepo := newMockOrgRepo() // Defined in organization_test.go
	svc := NewEmployeeService(empRepo, orgRepo, nil)

	ctx := context.Background()
	// Create org to pass existence check
	orgID, _ := orgRepo.CreateOrg(ctx, &model.Organization{Name: "Org1"})

	emp := &model.Employee{
		Name:           "John",
		Email:          "john@test.com",
		Phone:          "123",
		OrganizationID: orgID,
	}

	id, err := svc.Create(ctx, emp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 1 {
		t.Fatalf("expected id=1, got %d", id)
	}
}

func TestEmployeeCreateNoOrg(t *testing.T) {
	empRepo := newMockEmpRepo()
	orgRepo := newMockOrgRepo()
	svc := NewEmployeeService(empRepo, orgRepo, nil)

	emp := &model.Employee{
		Name:           "John",
		Email:          "john@test.com",
		Phone:          "123",
		OrganizationID: 999, // Does not exist
	}

	_, err := svc.Create(context.Background(), emp)
	if err == nil || err.Error() != "organization does not exist or is deleted" {
		t.Fatalf("expected missing org error, got %v", err)
	}
}

func TestEmployeeConflict(t *testing.T) {
	empRepo := newMockEmpRepo()
	orgRepo := newMockOrgRepo()
	svc := NewEmployeeService(empRepo, orgRepo, nil)
	ctx := context.Background()

	orgID, _ := orgRepo.CreateOrg(ctx, &model.Organization{Name: "Org1"})

	_, _ = svc.Create(ctx, &model.Employee{Name: "E1", Email: "dup@test.com", OrganizationID: orgID})

	_, err := svc.Create(ctx, &model.Employee{Name: "E2", Email: "dup@test.com", OrganizationID: orgID})
	if !errors.Is(err, repository.ErrConflict) {
		t.Fatalf("expected conflict on dup email, got %v", err)
	}
}
