package services

import (
	"nail_bot/internal/models"
	"nail_bot/internal/storage"
	"time"
)

type ReminderService struct{}

// GetBookingsForTomorrow возвращает записи на завтра
func (s *ReminderService) GetBookingsForTomorrow() ([]models.Booking, error) {
	var bookings []models.Booking
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	result := storage.GetDB().
		Where("date = ? AND status != ?", tomorrow, "cancelled").
		Find(&bookings)

	return bookings, result.Error
}
