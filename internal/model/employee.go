package model

type Employee struct {
	ID             int64  `db:"id"`
	Name           string `db:"name"`
	OrganizationID int64  `db:"organization_id"`
}
