package service

import (
	"QuickSlot/internal/model"
	"QuickSlot/internal/repository"
	"context"
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
