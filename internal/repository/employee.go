package repository

import (
	"QuickSlot/internal/model"
	"QuickSlot/pkg/database/mysql"
	"context"
	"errors"
)

var ErrEmployeeNotFound = errors.New("employee not found")

type EmployeeRepository interface {
	CreateEmployee(ctx context.Context, emp *model.Employee) (int64, error)
	GetByOrganization(ctx context.Context, orgID int64) ([]model.Employee, error)
	GetByID(ctx context.Context, id int64) (*model.Employee, error)
	Update(ctx context.Context, emp *model.Employee) error
	Delete(ctx context.Context, id int64) error
}

type employeeRepository struct {
	db *mysql.Dialect
}

func NewEmployeeRepository(db *mysql.Dialect) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) CreateEmployee(ctx context.Context, emp *model.Employee) (int64, error) {
	result, err := r.db.DB.ExecContext(ctx, "INSERT INTO employees (name, organization_id) VALUES (?, ?)",
		emp.Name, emp.OrganizationID,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (r *employeeRepository) GetByOrganization(ctx context.Context, orgID int64) ([]model.Employee, error) {
	employees := make([]model.Employee, 0)
	err := r.db.DB.SelectContext(ctx, &employees,
		"SELECT * FROM employees WHERE organization_id = ?", orgID)
	return employees, err
}

func (r *employeeRepository) GetByID(ctx context.Context, id int64) (*model.Employee, error) {
	var emp model.Employee
	err := r.db.DB.GetContext(ctx, &emp, "SELECT * FROM employees WHERE id = ?", id)
	if err != nil {
		return nil, ErrEmployeeNotFound
	}
	return &emp, nil
}

func (r *employeeRepository) Update(ctx context.Context, emp *model.Employee) error {
	_, err := r.db.DB.ExecContext(ctx, "UPDATE employees SET name = ? WHERE id = ?",
		emp.Name, emp.ID)
	return err
}

func (r *employeeRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.DB.ExecContext(ctx, "DELETE FROM employees WHERE id = ?", id)
	return err
}
