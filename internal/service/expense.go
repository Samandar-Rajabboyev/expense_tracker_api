package service

import (
	"context"
	"errors"
	"expense-tracker-api/internal/model"
	"expense-tracker-api/internal/repository"
	"time"
)

type ExpenseService struct {
	repo repository.ExpenseRepository
}

func NewExpenseService(repo *repository.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: *repo}
}

func (s *ExpenseService) Create(ctx context.Context, userID int64, title string, amount float64, category, description string, date string) (*model.Expense, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	expense := &model.Expense{
		Title:       title,
		Amount:      amount,
		Category:    category,
		Description: description,
		Date:        parsedDate,
		UserID:      userID,
	}

	if err := s.repo.Create(ctx, expense); err != nil {
		return nil, err
	}
	return expense, nil
}

func (s *ExpenseService) GetByUserID(ctx context.Context, userID int64) ([]*model.Expense, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *ExpenseService) GetByIDAndUserID(ctx context.Context, id int64, userID int64) (*model.Expense, error) {
	return s.repo.GetByIDAndUserID(ctx, id, userID)
}

func (s *ExpenseService) Update(ctx context.Context, expense *model.Expense) error {
	return s.repo.Update(ctx, expense)
}

func (s *ExpenseService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
