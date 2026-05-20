package repository

import (
	"context"
	"expense-tracker-api/internal/model"

	"gorm.io/gorm"
)

type ExpenseRepository struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) Create(ctx context.Context, expense *model.Expense) error {
	return r.db.WithContext(ctx).Create(expense).Error
}

func (r *ExpenseRepository) GetByID(ctx context.Context, id int64) (*model.Expense, error) {
	var expense model.Expense
	if err := r.db.WithContext(ctx).First(&expense, id).Error; err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *ExpenseRepository) GetAll(ctx context.Context) ([]*model.Expense, error) {
	var expenses []*model.Expense
	if err := r.db.WithContext(ctx).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *ExpenseRepository) GetByUserID(ctx context.Context, userID int64) ([]*model.Expense, error) {
	var expenses []*model.Expense
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *ExpenseRepository) GetByIDAndUserID(ctx context.Context, expenseID int64, userID int64) (*model.Expense, error) {
	var expense model.Expense
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", expenseID, userID).First(&expense).Error; err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *ExpenseRepository) GetByUserName(ctx context.Context, userName string) ([]*model.Expense, error) {
	var expenses []*model.Expense
	if err := r.db.WithContext(ctx).Joins("JOIN users ON users.id = expenses.user_id").Where("users.name = ?", userName).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *ExpenseRepository) Update(ctx context.Context, expense *model.Expense) error {
	return r.db.WithContext(ctx).Save(expense).Error
}

func (r *ExpenseRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Expense{}, id).Error
}
