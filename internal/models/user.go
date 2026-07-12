package models

import "time"

// Константы для шагов записи
const (
	StepService = "service"
	StepDate    = "date"
	StepTime    = "time"
	StepContact = "contact"
	StepConfirm = "confirm"
)

type User struct {
	ID         int64 `gorm:"primaryKey"`
	TelegramID int64 `gorm:"uniqueIndex;not null"`
	FirstName  string
	LastName   string
	Username   string
	Phone      string
	CreatedAt  time.Time
}

func (User) TableName() string {
	return "users"
}

type Booking struct {
	ID         uint   `gorm:"primaryKey"`
	UserID     int64  `gorm:"index;not null"`
	Service    string `gorm:"not null"`
	Date       string `gorm:"not null"`
	Time       string `gorm:"not null"`
	ClientName string `gorm:"not null"`
	Phone      string `gorm:"not null"`
	Status     string `gorm:"default:'pending'"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (Booking) TableName() string {
	return "bookings"
}

type UserSession struct {
	UserID    int64     `gorm:"primaryKey"`
	Step      string    `gorm:"not null"`
	Service   string    `gorm:"default:''"`
	Date      string    `gorm:"default:''"`
	Time      string    `gorm:"default:''"`
	Contact   string    `gorm:"default:''"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}
