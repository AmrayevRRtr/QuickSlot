package service

import (
	"QuickSlot/internal/model"
	"QuickSlot/internal/repository"
	"context"
	"errors"
	"testing"
)

type mockOrgRepo struct {
	orgs   map[int64]*model.Organization
	lastID int64
}

func newMockOrgRepo() *mockOrgRepo {
	return &mockOrgRepo{orgs: make(map[int64]*model.Organization)}
}

func (m *mockOrgRepo) CreateOrg(ctx context.Context, org *model.Organization) (int64, error) {
	m.lastID++
	org.ID = m.lastID
	m.orgs[m.lastID] = org
	return m.lastID, nil
}

func (m *mockOrgRepo) GetAll(ctx context.Context) ([]model.Organization, error) {
	var result []model.Organization
	for _, o := range m.orgs {
		result = append(result, *o)
	}
	return result, nil
}

func (m *mockOrgRepo) GetByID(ctx context.Context, id int64) (*model.Organization, error) {
	org, exists := m.orgs[id]
	if !exists {
		return nil, repository.ErrOrgNotFound
	}
	return org, nil
}

func (m *mockOrgRepo) ExistsByName(ctx context.Context, name string) (bool, error) {
	for _, o := range m.orgs {
		if o.Name == name {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockOrgRepo) Update(ctx context.Context, id int64, update *model.OrganizationUpdate) error {
	org, exists := m.orgs[id]
	if !exists {
		return repository.ErrOrgNotFound
	}
	if update.Name != nil {
		org.Name = *update.Name
	}
	return nil
}

func (m *mockOrgRepo) Delete(ctx context.Context, id int64) error {
	delete(m.orgs, id)
	return nil
}

func TestOrgCreateSuccess(t *testing.T) {
	repo := newMockOrgRepo()
	svc := NewOrganizationService(repo)

	id, err := svc.Create(context.Background(), "Google", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 1 {
		t.Fatalf("expected id=1, got %d", id)
	}
}

func TestOrgCreateConflict(t *testing.T) {
	repo := newMockOrgRepo()
	svc := NewOrganizationService(repo)

	_, _ = svc.Create(context.Background(), "Dup", 1)
	_, err := svc.Create(context.Background(), "Dup", 2)

	if !errors.Is(err, repository.ErrConflict) {
		t.Fatalf("expected conflict error, got %v", err)
	}
}

func TestOrgUpdate(t *testing.T) {
	repo := newMockOrgRepo()
	svc := NewOrganizationService(repo)

	id, _ := svc.Create(context.Background(), "Old", 1)

	newName := "New"
	err := svc.Update(context.Background(), id, &model.OrganizationUpdate{Name: &newName})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	org, _ := svc.GetByID(context.Background(), id)
	if org.Name != "New" {
		t.Fatalf("expected Name=New, got %s", org.Name)
	}
}
