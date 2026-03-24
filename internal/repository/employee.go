package repository

import (
	"QuickSlot/internal/model"
	"QuickSlot/pkg/database/mysql"
)

type EmployeeRepository interface {
	CreateEmployee(emp *model.Employee) (int64, error)
}

type employeeRepository struct {
	db *mysql.Dialect
}

func NewEmployeeRepository(db *mysql.Dialect) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) CreateEmployee(emp *model.Employee) (int64, error) {
	result, err := r.db.DB.Exec("INSERT INTO employees (name, organization_id) VALUES (?, ?)",
		emp.Name, emp.OrganizationID,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}
