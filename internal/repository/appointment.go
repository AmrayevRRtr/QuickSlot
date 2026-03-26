package repository

import (
	"QuickSlot/internal/model"
	"QuickSlot/pkg/database/mysql"
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrSlotAlreadyBooked   = errors.New("slot is already booked")
	ErrSlotNotFound        = errors.New("slot not found")
	ErrAppointmentNotFound = errors.New("appointment not found")
	ErrReviewNotFound      = errors.New("review not found")
)

type AppointmentRepository interface {
	BookSlot(ctx context.Context, userID, slotID int64) (*model.Appointment, error)
	CancelBooking(ctx context.Context, userID, appointmentID int64) error
	GetUserHistory(ctx context.Context, userID int64, from, to *time.Time) ([]model.AppointmentHistoryItem, error)
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

func (r *appointmentRepository) CancelBooking(ctx context.Context, userID, appointmentID int64) error {
	tx, err := r.db.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var slotID int64
	err = tx.GetContext(
		ctx,
		&slotID,
		"SELECT slot_id FROM appointments WHERE id=? AND user_id=? FOR UPDATE",
		appointmentID, userID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrAppointmentNotFound
		}
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE time_slots SET is_booked=FALSE WHERE id=?", slotID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM appointments WHERE id=? AND user_id=?", appointmentID, userID)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *appointmentRepository) GetUserHistory(ctx context.Context, userID int64, from, to *time.Time) ([]model.AppointmentHistoryItem, error) {
	items := make([]model.AppointmentHistoryItem, 0)

	query := `
SELECT
	a.id AS appointment_id,
	a.slot_id,
	ts.employee_id,
	ts.start_time,
	ts.end_time
FROM appointments a
JOIN time_slots ts ON ts.id = a.slot_id
WHERE a.user_id = ?
`
	args := []any{userID}

	if from != nil {
		query += " AND ts.start_time >= ?"
		args = append(args, *from)
	}
	if to != nil {
		query += " AND ts.start_time <= ?"
		args = append(args, *to)
	}

	query += " ORDER BY ts.start_time DESC"

	err := r.db.DB.SelectContext(ctx, &items, query, args...)
	return items, err
}
