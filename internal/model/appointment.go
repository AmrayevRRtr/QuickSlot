package model

type Appointment struct {
	ID     int64 `db:"id"`
	UserID int64 `db:"user_id"`
	SlotID int64 `db:"slot_id"`
}
