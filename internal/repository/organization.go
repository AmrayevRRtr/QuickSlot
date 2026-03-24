package repository

import (
	"QuickSlot/internal/model"
	"QuickSlot/pkg/database/mysql"
)

type OrganizationRepository interface {
	CreateOrg(org *model.Organization) (int64, error)
}

type organizationRepository struct {
	db *mysql.Dialect
}

func NewOrganizationRepository(db *mysql.Dialect) OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r organizationRepository) CreateOrg(org *model.Organization) (int64, error) {
	result, err := r.db.DB.Exec("INSERT INTO organizations (name, owner_id) VALUES (?, ?)",
		org.Name, org.OwnerID,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}
