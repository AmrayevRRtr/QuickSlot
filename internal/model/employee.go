package model

import "time"

type Employee struct {
	ID             int64      `db:"id" json:"id"`
	Name           string     `db:"name" json:"name"`
	Email          string     `db:"email" json:"email"`
	Phone          string     `db:"phone" json:"phone"`
	OrganizationID int64      `db:"organization_id" json:"organization_id"`
	DeletedAt      *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type EmployeeUpdate struct {
	Name           *string `db:"name" json:"name,omitempty"`
	Email          *string `db:"email" json:"email,omitempty"`
	Phone          *string `db:"phone" json:"phone,omitempty"`
	OrganizationID *int64  `db:"organization_id" json:"organization_id,omitempty"`
}
