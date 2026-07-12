package services

import (
	"errors"
	"nail_bot/internal/models"
	"nail_bot/internal/storage"
	"time"
)

var MasterName = "Инесса"

type BookingService struct{}

func (s *BookingService) CreateBooking(userID int64, service, date, timeStr, clientName, phone string) (*models.Booking, error) {
	// Проверяем занятость
	exists, err := s.CheckAvailability(date, timeStr)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("это время уже занято")
	}

	booking := &models.Booking{
		UserID:     userID,
		Service:    service,
		Date:       date,
		Time:       timeStr,
		ClientName: clientName,
		Phone:      phone,
		Status:     "pending",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	result := storage.GetDB().Create(booking)
	if result.Error != nil {
		return nil, result.Error
	}

	return booking, nil
}

func (s *BookingService) CheckAvailability(date, timeStr string) (bool, error) {
	var count int64
	result := storage.GetDB().Model(&models.Booking{}).
		Where("date = ? AND time = ? AND status != ?", date, timeStr, "cancelled").
		Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

func (s *BookingService) GetUserBookings(userID int64) ([]models.Booking, error) {
	var bookings []models.Booking
	result := storage.GetDB().
		Where("user_id = ? AND status != ?", userID, "cancelled").
		Order("date ASC, time ASC").
		Find(&bookings)

	return bookings, result.Error
}

func (s *BookingService) GetAllBookings() ([]models.Booking, error) {
	var bookings []models.Booking
	result := storage.GetDB().
		Where("status != ?", "cancelled").
		Order("date ASC, time ASC").
		Find(&bookings)

	return bookings, result.Error
}

func (s *BookingService) GetBookingsForNext3Days() ([]models.Booking, error) {
	var bookings []models.Booking
	today := time.Now().Format("2006-01-02")
	threeDaysLater := time.Now().AddDate(0, 0, 3).Format("2006-01-02")

	result := storage.GetDB().
		Where("date >= ? AND date <= ? AND status != ?", today, threeDaysLater, "cancelled").
		Order("date ASC, time ASC").
		Find(&bookings)

	return bookings, result.Error
}

func (s *BookingService) CancelBooking(bookingID uint, userID int64) error {
	result := storage.GetDB().
		Model(&models.Booking{}).
		Where("id = ? AND user_id = ?", bookingID, userID).
		Update("status", "cancelled")

	if result.RowsAffected == 0 {
		return errors.New("запись не найдена")
	}
	return result.Error
}

// 12 июля
// GetActiveBookingsByUserID возвращает активные записи пользователя
func (s *BookingService) GetActiveBookingsByUserID(userID int64) ([]models.Booking, error) {
	var bookings []models.Booking
	result := storage.GetDB().
		Where("user_id = ? AND status NOT IN (?)", userID, []string{"cancelled", "completed"}).
		Order("date ASC, time ASC").
		Find(&bookings)

	if result.Error != nil {
		return nil, result.Error
	}
	return bookings, nil
}

// CancelBookingByID отменяет запись по ID и userID
func (s *BookingService) CancelBookingByID(bookingID uint, userID int64) error {
	result := storage.GetDB().
		Model(&models.Booking{}).
		Where("id = ? AND user_id = ?", bookingID, userID).
		Update("status", "cancelled")

	if result.RowsAffected == 0 {
		return errors.New("запись не найдена или уже отменена")
	}
	return result.Error
}
