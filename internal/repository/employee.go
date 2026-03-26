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
	ErrEmployeeNotFound = errors.New("employee not found")
)

type EmployeeRepository interface {
	CreateEmployee(ctx context.Context, emp *model.Employee) (int64, error)
	GetByOrganization(ctx context.Context, orgID int64) ([]model.Employee, error)
	GetByID(ctx context.Context, id int64) (*model.Employee, error)
	ExistsByEmailOrPhone(ctx context.Context, email, phone string) (bool, error)
	Update(ctx context.Context, id int64, update *model.EmployeeUpdate) error
	Delete(ctx context.Context, id int64) error
}

type employeeRepository struct {
	db *mysql.Dialect
}

func NewEmployeeRepository(db *mysql.Dialect) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) ExistsByEmailOrPhone(ctx context.Context, email, phone string) (bool, error) {
	var exists bool
	query := `
        SELECT EXISTS(
            SELECT 1 FROM employees 
            WHERE (email = ? OR phone = ?) AND deleted_at IS NULL
        )
    `
	err := r.db.DB.GetContext(ctx, &exists, query, email, phone)
	return exists, err
}

func (r *employeeRepository) CreateEmployee(ctx context.Context, emp *model.Employee) (int64, error) {
	result, err := r.db.DB.ExecContext(ctx, "INSERT INTO employees (name, email, phone, organization_id) VALUES (?, ?, ?, ?)",
		emp.Name, emp.Email, emp.Phone, emp.OrganizationID,
	)
	if err != nil {
		var mysqlErr *mysqlDriver.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return 0, ErrConflict
		}
		return 0, err
	}

	return result.LastInsertId()
}

func (r *employeeRepository) GetByOrganization(ctx context.Context, orgID int64) ([]model.Employee, error) {
	employees := make([]model.Employee, 0)
	err := r.db.DB.SelectContext(ctx, &employees,
		"SELECT * FROM employees WHERE organization_id = ? AND deleted_at IS NULL", orgID)
	return employees, err
}

func (r *employeeRepository) GetByID(ctx context.Context, id int64) (*model.Employee, error) {
	var emp model.Employee
	err := r.db.DB.GetContext(ctx, &emp, "SELECT * FROM employees WHERE id = ? AND deleted_at IS NULL", id)
	if err != nil {
		return nil, ErrEmployeeNotFound
	}
	return &emp, nil
}

func (r *employeeRepository) Update(ctx context.Context, id int64, update *model.EmployeeUpdate) error {
	var setClauses []string
	var args []interface{}

	if update.Name != nil {
		setClauses = append(setClauses, "name = ?")
		args = append(args, *update.Name)
	}
	if update.Email != nil {
		setClauses = append(setClauses, "email = ?")
		args = append(args, *update.Email)
	}
	if update.Phone != nil {
		setClauses = append(setClauses, "phone = ?")
		args = append(args, *update.Phone)
	}
	if update.OrganizationID != nil {
		setClauses = append(setClauses, "organization_id = ?")
		args = append(args, *update.OrganizationID)
	}

	if len(setClauses) == 0 {
		return nil
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE employees SET %s WHERE id = ? AND deleted_at IS NULL", strings.Join(setClauses, ", "))

	_, err := r.db.DB.ExecContext(ctx, query, args...)
	if err != nil {
		var mysqlErr *mysqlDriver.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return ErrConflict
		}
		return err
	}
	return nil
}

func (r *employeeRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.DB.ExecContext(ctx, "UPDATE employees SET deleted_at = NOW() WHERE id = ?", id)
	return err
}
