package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"QuickSlot/internal/model"
	"QuickSlot/pkg/database/mysql"
	mysqlDriver "github.com/go-sql-driver/mysql"
)

var (
	ErrOrgNotFound = errors.New("organization not found")
	ErrConflict    = errors.New("conflict: record already exists")
)

type OrganizationRepository interface {
	CreateOrg(ctx context.Context, org *model.Organization) (int64, error)
	GetAll(ctx context.Context) ([]model.Organization, error)
	GetByID(ctx context.Context, id int64) (*model.Organization, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
	Update(ctx context.Context, id int64, update *model.OrganizationUpdate) error
	Delete(ctx context.Context, id int64) error
}

type organizationRepository struct {
	db *mysql.Dialect
}

func NewOrganizationRepository(db *mysql.Dialect) OrganizationRepository {
	return &organizationRepository{db: db}
}

func isMySQLDuplicateError(err error) bool {
	var mysqlErr *mysqlDriver.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return true
	}
	return false
}

func (r *organizationRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var exists bool
	err := r.db.DB.GetContext(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM organizations WHERE name = ? AND deleted_at IS NULL)", name)
	return exists, err
}

func (r *organizationRepository) CreateOrg(ctx context.Context, org *model.Organization) (int64, error) {
	result, err := r.db.DB.ExecContext(ctx, "INSERT INTO organizations (name, owner_id) VALUES (?, ?)",
		org.Name, org.OwnerID,
	)
	if err != nil {
		if isMySQLDuplicateError(err) {
			return 0, ErrConflict
		}
		return 0, err
	}
	return result.LastInsertId()
}

func (r *organizationRepository) GetAll(ctx context.Context) ([]model.Organization, error) {
	orgs := make([]model.Organization, 0)
	err := r.db.DB.SelectContext(ctx, &orgs, "SELECT * FROM organizations WHERE deleted_at IS NULL")
	return orgs, err
}

func (r *organizationRepository) GetByID(ctx context.Context, id int64) (*model.Organization, error) {
	var org model.Organization
	err := r.db.DB.GetContext(ctx, &org, "SELECT * FROM organizations WHERE id = ? AND deleted_at IS NULL", id)
	if err != nil {
		return nil, ErrOrgNotFound
	}
	return &org, nil
}

func (r *organizationRepository) Update(ctx context.Context, id int64, update *model.OrganizationUpdate) error {
	var setClauses []string
	var args []interface{}

	if update.Name != nil {
		setClauses = append(setClauses, "name = ?")
		args = append(args, *update.Name)
	}

	if len(setClauses) == 0 {
		return nil // Nothing to update
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE organizations SET %s WHERE id = ? AND deleted_at IS NULL", strings.Join(setClauses, ", "))

	_, err := r.db.DB.ExecContext(ctx, query, args...)
	if err != nil {
		if isMySQLDuplicateError(err) {
			return ErrConflict
		}
		return err
	}
	return nil
}

func (r *organizationRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.DB.ExecContext(ctx, "UPDATE organizations SET deleted_at = NOW() WHERE id = ?", id)
	return err
}
