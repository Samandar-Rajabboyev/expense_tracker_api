package service

import (
	"context"
	"expense-tracker-api/internal/model"
	"testing"
)

type FakeExpenseRepository struct{}

func (f *FakeExpenseRepository) Create(ctx context.Context, expense *model.Expense) error {
	return nil
}
func (f *FakeExpenseRepository) GetByID(ctx context.Context, id int64) (*model.Expense, error) {
	return nil, nil
}
func (f *FakeExpenseRepository) GetAll(ctx context.Context) ([]*model.Expense, error) {
	return nil, nil
}
func (f *FakeExpenseRepository) GetByUserID(ctx context.Context, userID int64) ([]*model.Expense, error) {
	return nil, nil
}
func (f *FakeExpenseRepository) GetByIDAndUserID(ctx context.Context, expenseID int64, userID int64) (*model.Expense, error) {
	return nil, nil
}
func (f *FakeExpenseRepository) GetByUserName(ctx context.Context, userName string) ([]*model.Expense, error) {
	return nil, nil
}
func (f *FakeExpenseRepository) Update(ctx context.Context, expense *model.Expense) error {
	return nil
}
func (f *FakeExpenseRepository) Delete(ctx context.Context, id int64) error {
	return nil
}

func TestCreateExpense_Validation(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		amount  float64
		date    string
		wantErr bool
	}{
		{"empty title", "", 100, "2026-01-01", true},
		{"negative amount", "Lunch", -50, "2026-01-01", true},
		{"invalid date", "Lunch", 50, "not-a-date", true},
		{"valid input", "Lunch", 50, "2026-01-01", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &FakeExpenseRepository{}
			svc := NewExpenseService(repo)
			_, err := svc.Create(context.Background(), 1, tt.title, tt.amount, "food", "lunch", tt.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
			}
		})
	}
}
