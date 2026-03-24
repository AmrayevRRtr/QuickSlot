package repository

import (
	"QuickSlot/internal/model"
	"QuickSlot/pkg/database/mysql"
	"context"
	"database/sql"
	"errors"
)

var (
	ErrSlotAlreadyBooked = errors.New("slot is already booked")
	ErrSlotNotFound      = errors.New("slot not found")
)

type AppointmentRepository interface {
	BookSlot(ctx context.Context, userID, slotID int64) (*model.Appointment, error)
}

type appointmentRepository struct {
	db *mysql.Dialect
}

func NewAppointmentRepository(db *mysql.Dialect) AppointmentRepository {
	return &appointmentRepository{db: db}
}

func (r *appointmentRepository) BookSlot(ctx context.Context, userID, slotID int64) (*model.Appointment, error) {
	tx, err := r.db.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	var isBooked bool
	err = tx.GetContext(ctx, &isBooked,
		"SELECT is_booked FROM time_slots WHERE id=? FOR UPDATE", slotID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrSlotNotFound
		}
		return nil, err
	}

	if isBooked {
		return nil, ErrSlotAlreadyBooked
	}

	_, err = tx.ExecContext(ctx, "UPDATE time_slots SET is_booked=TRUE WHERE id=?", slotID)
	if err != nil {
		return nil, err
	}

	res, err := tx.ExecContext(ctx,
		"INSERT INTO appointments (user_id, slot_id) VALUES (?, ?)",
		userID, slotID,
	)

	if err != nil {
		return nil, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	var appointment model.Appointment
	err = tx.GetContext(ctx, &appointment, "SELECT * FROM appointments WHERE id=?", lastID)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &appointment, nil
}
