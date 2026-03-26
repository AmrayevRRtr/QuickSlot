package model

import "time"

type Appointment struct {
	ID     int64 `db:"id"`
	UserID int64 `db:"user_id"`
	SlotID int64 `db:"slot_id"`
}

type AppointmentHistoryItem struct {
	AppointmentID int64 `db:"appointment_id"`
	SlotID        int64 `db:"slot_id"`
	EmployeeID    int64 `db:"employee_id"`
	StartTime     time.Time `db:"start_time"`
	EndTime       time.Time `db:"end_time"`
}
