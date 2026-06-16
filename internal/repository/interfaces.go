package repository

import (
	"context"
	"expense-tracker-api/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByIDWithExpenses(ctx context.Context, id int64) (*model.User, error)
	CreateUserWithExpense(ctx context.Context, user *model.User, expense *model.Expense) error
	DeleteUserWithExpense(ctx context.Context, id int64) error
}

type ExpenseRepository interface {
	Create(ctx context.Context, expense *model.Expense) error
	GetByID(ctx context.Context, id int64) (*model.Expense, error)
	GetAll(ctx context.Context) ([]*model.Expense, error)
	GetByUserID(ctx context.Context, userID int64) ([]*model.Expense, error)
	GetByIDAndUserID(ctx context.Context, expenseID int64, userID int64) (*model.Expense, error)
	GetByUserName(ctx context.Context, userName string) ([]*model.Expense, error)
	Update(ctx context.Context, expense *model.Expense) error
	Delete(ctx context.Context, id int64) error
}
