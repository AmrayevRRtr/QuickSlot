package model

import "time"

type Review struct {
	ID             int64     `db:"id" json:"id"`
	UserID         int64     `db:"user_id" json:"user_id"`
	OrganizationID int64     `db:"organization_id" json:"organization_id"`
	Rating         int       `db:"rating" json:"rating"`
	Comment        string    `db:"comment" json:"comment"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}
