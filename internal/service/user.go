package service

import (
	"context"
	"errors"
	"expense-tracker-api/internal/model"
	"expense-tracker-api/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, email, password, name string) (*model.User, error) {
	existing, _ := s.repo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &model.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, err
}

func (s *UserService) Login(ctx context.Context, email, password string) (*model.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	return s.repo.GetByIDWithExpenses(ctx, id)
}
