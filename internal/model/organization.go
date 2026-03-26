package model

import "time"

type Organization struct {
	ID        int64      `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	OwnerID   int64      `db:"owner_id" json:"owner_id"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type OrganizationUpdate struct {
	Name *string `db:"name" json:"name,omitempty"`
}
