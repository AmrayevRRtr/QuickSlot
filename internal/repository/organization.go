package repository

import (
	"QuickSlot/internal/model"
	"QuickSlot/pkg/database/mysql"
	"context"
	"errors"
)

var ErrOrgNotFound = errors.New("organization not found")

type OrganizationRepository interface {
	CreateOrg(ctx context.Context, org *model.Organization) (int64, error)
	GetAll(ctx context.Context) ([]model.Organization, error)
	GetByID(ctx context.Context, id int64) (*model.Organization, error)
	Update(ctx context.Context, org *model.Organization) error
	Delete(ctx context.Context, id int64) error
}

type organizationRepository struct {
	db *mysql.Dialect
}

func NewOrganizationRepository(db *mysql.Dialect) OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) CreateOrg(ctx context.Context, org *model.Organization) (int64, error) {
	result, err := r.db.DB.ExecContext(ctx, "INSERT INTO organizations (name, owner_id) VALUES (?, ?)",
		org.Name, org.OwnerID,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (r *organizationRepository) GetAll(ctx context.Context) ([]model.Organization, error) {
	orgs := make([]model.Organization, 0)
	err := r.db.DB.SelectContext(ctx, &orgs, "SELECT * FROM organizations")
	return orgs, err
}

func (r *organizationRepository) GetByID(ctx context.Context, id int64) (*model.Organization, error) {
	var org model.Organization
	err := r.db.DB.GetContext(ctx, &org, "SELECT * FROM organizations WHERE id = ?", id)
	if err != nil {
		return nil, ErrOrgNotFound
	}
	return &org, nil
}

func (r *organizationRepository) Update(ctx context.Context, org *model.Organization) error {
	_, err := r.db.DB.ExecContext(ctx, "UPDATE organizations SET name = ? WHERE id = ?",
		org.Name, org.ID)
	return err
}

func (r *organizationRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.DB.ExecContext(ctx, "DELETE FROM organizations WHERE id = ?", id)
	return err
}
