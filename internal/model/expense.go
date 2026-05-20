package model

import "time"

type Expense struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string    `gorm:"size:100;not null" json:"title"`
	Amount      float64   `gorm:"type:numeric(10,2);not null" json:"amount"`
	Category    string    `gorm:"size:50" json:"category"`
	Description string    `gorm:"type:text" json:"description"`
	Date        time.Time `gorm:"column:expense_date;type:date;not null" json:"date"`
	UserID      int64     `gorm:"index;not null" json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	User        *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
