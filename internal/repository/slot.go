package repository

import (
	"QuickSlot/internal/model"
	"QuickSlot/pkg/database/mysql"
)

type SlotRepository interface {
	BulkCreate(slots []model.Slot) error
	GetAvailableByEmployee(employeeID int64) ([]model.Slot, error)
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

func (r *slotRepository) GetAvailableByEmployee(employeeID int64) ([]model.Slot, error) {
	var slots []model.Slot

	err := r.db.DB.Select(&slots,
		`SELECT * FROM time_slots WHERE employee_id=? AND is_booked=FALSE ORDER BY start_time`,
		employeeID,
	)

	return slots, err
}
