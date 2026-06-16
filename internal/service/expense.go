package service

import (
	"context"
	"encoding/json"
	"errors"
	"expense-tracker-api/internal/cache"
	"expense-tracker-api/internal/model"
	"expense-tracker-api/internal/repository"
	"fmt"
	"time"
)

type ExpenseService struct {
	repo  repository.ExpenseRepository
	cache *cache.Cache
}

func NewExpenseService(repo repository.ExpenseRepository, cache *cache.Cache) *ExpenseService {
	return &ExpenseService{repo: repo, cache: cache}
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

	cacheKey := fmt.Sprintf("user:%d:expenses", userID)
	s.cache.Delete(ctx, cacheKey)

	return expense, nil
}

func (s *ExpenseService) GetByUserID(ctx context.Context, userID int64) ([]*model.Expense, error) {
	cacheKey := fmt.Sprintf("user:%d:expenses", userID)

	cached, err := s.cache.Get(ctx, cacheKey)
	if err == nil {
		var expenses []*model.Expense
		if jsonErr := json.Unmarshal([]byte(cached), &expenses); jsonErr == nil {
			return expenses, nil
		}
	}

	expenses, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(expenses)
	s.cache.Set(ctx, cacheKey, string(data), 5*time.Minute)

	return expenses, nil
}

func (s *ExpenseService) GetByIDAndUserID(ctx context.Context, id int64, userID int64) (*model.Expense, error) {
	return s.repo.GetByIDAndUserID(ctx, id, userID)
}

func (s *ExpenseService) Update(ctx context.Context, expense *model.Expense) error {
	existing, err := s.repo.GetByID(ctx, expense.ID)
	if err != nil {
		return err
	}

	err = s.repo.Update(ctx, expense)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("user:%d:expenses", existing.UserID)
	s.cache.Delete(ctx, cacheKey)

	return nil
}

func (s *ExpenseService) Delete(ctx context.Context, id int64) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("user:%d:expenses", existing.UserID)
	s.cache.Delete(ctx, cacheKey)

	return nil
}
