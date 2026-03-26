package repository

import (
	"QuickSlot/internal/model"
	"QuickSlot/pkg/database/mysql"
	"time"
)

type SlotRepository interface {
	BulkCreate(slots []model.Slot) error
	GetAvailableByEmployee(employeeID int64, from, to *time.Time) ([]model.Slot, error)
}

type slotRepository struct {
	db *mysql.Dialect
}

func NewSlotRepository(db *mysql.Dialect) SlotRepository {
	return &slotRepository{db: db}
}

func (r *slotRepository) BulkCreate(slots []model.Slot) error {
	query := `INSERT IGNORE INTO time_slots (employee_id, start_time, end_time) VALUES (?, ?, ?)`

	for _, slot := range slots {
		_, err := r.db.DB.Exec(query,
			slot.EmployeeID, slot.StartTime, slot.EndTime)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *slotRepository) GetAvailableByEmployee(employeeID int64, from, to *time.Time) ([]model.Slot, error) {
	var slots []model.Slot

	query := `
SELECT *
FROM time_slots
WHERE employee_id = ?
  AND is_booked = FALSE
`
	args := []any{employeeID}
	if from != nil {
		query += " AND start_time >= ?"
		args = append(args, *from)
	}
	if to != nil {
		query += " AND start_time <= ?"
		args = append(args, *to)
	}

	query += " ORDER BY start_time"

	err := r.db.DB.Select(&slots, query, args...)

	return slots, err
}
