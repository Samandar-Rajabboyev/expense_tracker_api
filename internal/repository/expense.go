package repository

import (
	"context"
	"expense-tracker-api/internal/model"

	"gorm.io/gorm"
)

type ExpenseRepositoryImpl struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) *ExpenseRepositoryImpl {
	return &ExpenseRepositoryImpl{db: db}
}

func (r *ExpenseRepositoryImpl) Create(ctx context.Context, expense *model.Expense) error {
	return r.db.WithContext(ctx).Create(expense).Error
}

func (r *ExpenseRepositoryImpl) GetByID(ctx context.Context, id int64) (*model.Expense, error) {
	var expense model.Expense
	if err := r.db.WithContext(ctx).First(&expense, id).Error; err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *ExpenseRepositoryImpl) GetAll(ctx context.Context) ([]*model.Expense, error) {
	var expenses []*model.Expense
	if err := r.db.WithContext(ctx).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *ExpenseRepositoryImpl) GetByUserID(ctx context.Context, userID int64) ([]*model.Expense, error) {
	var expenses []*model.Expense
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *ExpenseRepositoryImpl) GetByIDAndUserID(ctx context.Context, expenseID int64, userID int64) (*model.Expense, error) {
	var expense model.Expense
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", expenseID, userID).First(&expense).Error; err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *ExpenseRepositoryImpl) GetByUserName(ctx context.Context, userName string) ([]*model.Expense, error) {
	var expenses []*model.Expense
	if err := r.db.WithContext(ctx).Joins("JOIN users ON users.id = expenses.user_id").Where("users.name = ?", userName).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *ExpenseRepositoryImpl) Update(ctx context.Context, expense *model.Expense) error {
	return r.db.WithContext(ctx).Save(expense).Error
}

func (r *ExpenseRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Expense{}, id).Error
}
