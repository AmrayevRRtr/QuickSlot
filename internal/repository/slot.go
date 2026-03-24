package repository

import (
	"QuickSlot/internal/model"
	"QuickSlot/pkg/database/mysql"
)

type SlotRepository interface {
	CreateSlot(slot *model.Slot) error
	GetAvailableByEmployee(employeeID int64) ([]model.Slot, error)
}

type slotRepository struct {
	db *mysql.Dialect
}

func NewSlotRepository(db *mysql.Dialect) SlotRepository {
	return &slotRepository{db: db}
}

func (r *slotRepository) CreateSlot(slot *model.Slot) error {
	_, err := r.db.DB.Exec(
		`INSERT INTO time_slots (employee_id, start_time, end_time) VALUES (?, ?, ?)`,
		slot.EmployeeID, slot.StartTime, slot.EndTime,
	)
	return err
}

func (r *slotRepository) GetAvailableByEmployee(employeeID int64) ([]model.Slot, error) {
	var slots []model.Slot

	err := r.db.DB.Select(&slots,
		`SELECT * FROM time_slots WHERE employee_id=? AND is_booked=FALSE`,
		employeeID,
	)

	return slots, err
}
