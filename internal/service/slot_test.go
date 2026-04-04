package service

import (
	"QuickSlot/internal/model"
	"context"
	"testing"
	"time"
)

type mockSlotRepo struct {
	slots []model.Slot
}

func (m *mockSlotRepo) BulkCreate(slots []model.Slot) error {
	m.slots = append(m.slots, slots...)
	return nil
}

func (m *mockSlotRepo) GetAvailableByEmployee(employeeID int64, from, to *time.Time) ([]model.Slot, error) {
	var result []model.Slot
	for _, s := range m.slots {
		if s.EmployeeID == employeeID && !s.IsBooked {
			result = append(result, s)
		}
	}
	return result, nil
}

func TestGenerateSlots(t *testing.T) {
	repo := &mockSlotRepo{}
	svc := NewSlotService(repo)

	date := time.Date(2026, 3, 25, 0, 0, 0, 0, time.UTC)
	err := svc.GenerateSlots(context.Background(), 1, date, 9, 12, 30)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// 9:00-12:00 with 30min slots = 6 slots
	if len(repo.slots) != 6 {
		t.Fatalf("expected 6 slots, got %d", len(repo.slots))
	}

	// check first slot
	first := repo.slots[0]
	if first.StartTime.Hour() != 9 || first.StartTime.Minute() != 0 {
		t.Fatalf("first slot should start at 9:00, got %v", first.StartTime)
	}
	if first.EndTime.Hour() != 9 || first.EndTime.Minute() != 30 {
		t.Fatalf("first slot should end at 9:30, got %v", first.EndTime)
	}

	// check last slot
	last := repo.slots[5]
	if last.StartTime.Hour() != 11 || last.StartTime.Minute() != 30 {
		t.Fatalf("last slot should start at 11:30, got %v", last.StartTime)
	}
}

func TestGenerateSlotsEmployeeID(t *testing.T) {
	repo := &mockSlotRepo{}
	svc := NewSlotService(repo)

	date := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	_ = svc.GenerateSlots(context.Background(), 42, date, 10, 11, 30)

	for _, s := range repo.slots {
		if s.EmployeeID != 42 {
			t.Fatalf("expected employee_id=42, got %d", s.EmployeeID)
		}
	}
}

func TestGenerateSlotsEmpty(t *testing.T) {
	repo := &mockSlotRepo{}
	svc := NewSlotService(repo)

	date := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	err := svc.GenerateSlots(context.Background(), 1, date, 10, 10, 30)
	if err != nil {
		t.Fatal(err)
	}
	if len(repo.slots) != 0 {
		t.Fatalf("expected 0 slots when start=end, got %d", len(repo.slots))
	}
}
