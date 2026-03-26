package service

import (
	"QuickSlot/internal/model"
	"QuickSlot/internal/repository"
	"context"
	"time"
)

type SlotService struct {
	repo repository.SlotRepository
}

func NewSlotService(r repository.SlotRepository) *SlotService {
	return &SlotService{repo: r}
}

func (s *SlotService) GenerateSlots(
	ctx context.Context,
	employeeID int64,
	date time.Time,
	startHour int,
	endHour int,
	durationMinutes int,
) error {

	start := time.Date(date.Year(), date.Month(), date.Day(), startHour, 0, 0, 0, time.UTC)
	end := time.Date(date.Year(), date.Month(), date.Day(), endHour, 0, 0, 0, time.UTC)

	var slots []model.Slot

	for t := start; t.Before(end); t = t.Add(time.Duration(durationMinutes) * time.Minute) {
		slot := model.Slot{
			EmployeeID: employeeID,
			StartTime:  t,
			EndTime:    t.Add(time.Duration(durationMinutes) * time.Minute),
		}
		slots = append(slots, slot)
	}

	return s.repo.BulkCreate(slots)
}

func (s *SlotService) GetAvailableByEmployee(ctx context.Context, employeeID int64, from, to *time.Time) ([]model.Slot, error) {
	return s.repo.GetAvailableByEmployee(employeeID, from, to)
}
