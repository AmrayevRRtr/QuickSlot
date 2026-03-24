package model

type Organization struct {
	ID      int64  `db:"id"`
	Name    string `db:"name"`
	OwnerID int64  `db:"owner_id"`
}
