package service

import (
	"context"
	"errors"
	"expense-tracker-api/internal/model"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type FakeUserRepository struct {
	GetByEmailFunc func(ctx context.Context, email string) (*model.User, error)
	CreateFunc     func(ctx context.Context, user *model.User) error
}

func (f *FakeUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return f.GetByEmailFunc(ctx, email)
}
func (f *FakeUserRepository) Create(ctx context.Context, user *model.User) error {
	return f.CreateFunc(ctx, user)
}

func (f *FakeUserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	return nil, nil
}
func (f *FakeUserRepository) GetByIDWithExpenses(ctx context.Context, id int64) (*model.User, error) {
	return nil, nil
}
func (f *FakeUserRepository) CreateUserWithExpense(ctx context.Context, user *model.User, expense *model.Expense) error {
	return nil
}
func (f *FakeUserRepository) DeleteUserWithExpense(ctx context.Context, id int64) error {
	return nil
}

func TestUserService_Register(t *testing.T) {
	tests := []struct {
		name           string
		getByEmailFunc func(ctx context.Context, email string) (*model.User, error)
		createFunc     func(ctx context.Context, user *model.User) error
		wantErr        bool
	}{
		{
			name: "email already registered",
			getByEmailFunc: func(ctx context.Context, email string) (*model.User, error) {
				return &model.User{Email: email}, nil // user exists
			},
			createFunc: nil, // won't be called
			wantErr:    true,
		},
		{
			name: "successful registration",
			getByEmailFunc: func(ctx context.Context, email string) (*model.User, error) {
				return nil, errors.New("not found") // no existing user
			},
			createFunc: func(ctx context.Context, user *model.User) error {
				return nil // create succeeds
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &FakeUserRepository{
				GetByEmailFunc: tt.getByEmailFunc,
				CreateFunc:     tt.createFunc,
			}
			svc := NewUserService(fake)
			_, err := svc.Register(context.Background(), "test@test.com", "password123", "Test User")
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
	hashedPw, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	tests := []struct {
		name           string
		getByEmailFunc func(ctx context.Context, email string) (*model.User, error)
		password       string
		wantErr        bool
	}{
		{
			name: "user not found",
			getByEmailFunc: func(ctx context.Context, email string) (*model.User, error) {
				return nil, errors.New("not found")
			},
			password: "wrong",
			wantErr:  true,
		},
		{
			name: "wrong password",
			getByEmailFunc: func(ctx context.Context, email string) (*model.User, error) {
				return &model.User{Email: email, Password: string(hashedPw)}, nil
			},
			password: "wrong",
			wantErr:  true,
		},
		{
			name: "successful login",
			getByEmailFunc: func(ctx context.Context, email string) (*model.User, error) {
				return &model.User{Email: email, Password: string(hashedPw)}, nil
			},
			password: "correct",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &FakeUserRepository{
				GetByEmailFunc: tt.getByEmailFunc,
			}
			svc := NewUserService(fake)
			_, err := svc.Login(context.Background(), "test@test.com", tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
			}
		})
	}
}
