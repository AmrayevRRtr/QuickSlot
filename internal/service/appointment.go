package service

import (
	"QuickSlot/internal/model"
	"QuickSlot/internal/repository"
	"context"
	"time"
)

type AppointmentService struct {
	repo repository.AppointmentRepository
}

func NewAppointmentService(r repository.AppointmentRepository) *AppointmentService {
	return &AppointmentService{repo: r}
}

func (s *AppointmentService) BookSlot(ctx context.Context, userID, slotID int64) (*model.Appointment, error) {
	return s.repo.BookSlot(ctx, userID, slotID)
}

func (s *AppointmentService) CancelBooking(ctx context.Context, userID, appointmentID int64) error {
	return s.repo.CancelBooking(ctx, userID, appointmentID)
}

func (s *AppointmentService) GetUserHistory(ctx context.Context, userID int64, from, to *time.Time) ([]model.AppointmentHistoryItem, error) {
	return s.repo.GetUserHistory(ctx, userID, from, to)
}
