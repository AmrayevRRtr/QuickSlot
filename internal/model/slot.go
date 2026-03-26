package model

import "time"

type Slot struct {
	ID         int64     `db:"id"`
	EmployeeID int64     `db:"employee_id"`
	StartTime  time.Time `db:"start_time"`
	EndTime    time.Time `db:"end_time"`
	IsBooked   bool      `db:"is_booked"`
}
