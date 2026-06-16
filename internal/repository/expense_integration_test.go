package repository

import (
	"context"
	"expense-tracker-api/internal/config"
	"expense-tracker-api/internal/database"
	"expense-tracker-api/internal/model"
	// "os"
	"testing"
	"time"

	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	config, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load configurations")
	}
	dbURL := config.DatabaseURL
	if dbURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}
	db, err := database.NewDB(dbURL)
	if err != nil {
		t.Fatalf("failed to connect to test DB: %v", err)
	}
	// Auto-migrate the test tables
	db.AutoMigrate(&model.User{}, &model.Expense{})
	return db
}
func TestExpenseRepository_CreateAndGet(t *testing.T) {
	db := setupTestDB(t)
	repo := NewExpenseRepository(db)
	// First create a user (expenses need a user_id)
	userRepo := NewUserRepository(db)
	user := &model.User{
		Email:    "integration@test.com",
		Password: "hashed",
		Name:     "Integration Test",
	}
	userRepo.Create(context.Background(), user)
	// Now test expense Create
	expense := &model.Expense{
		Title:       "Test Expense",
		Amount:      50.0,
		Category:    "food",
		Description: "Test lunch",
		Date:        time.Now(),
		UserID:      user.ID,
	}
	err := repo.Create(context.Background(), expense)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	// Test GetByID
	got, err := repo.GetByID(context.Background(), expense.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got.Title != "Test Expense" {
		t.Errorf("expected title 'Test Expense', got %q", got.Title)
	}
	if got.Amount != 50.0 {
		t.Errorf("expected amount 50.0, got %v", got.Amount)
	}
}
