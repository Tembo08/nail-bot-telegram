package services

import (
	"regexp"
	"testing"
	"time"

	"nail_bot/internal/models"
	"nail_bot/internal/storage"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupTestDBForReport создаёт мок-базу для тестов отчётов
func setupTestDBForReport(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open gorm: %v", err)
	}

	storage.DB = gormDB
	return gormDB, mock
}

// Тест для GetBookingsForNext3Days — есть записи
func TestGetBookingsForNext3Days_WithBookings(t *testing.T) {
	_, mock := setupTestDBForReport(t)

	today := time.Now().Format("2006-01-02")
	threeDaysLater := time.Now().AddDate(0, 0, 3).Format("2006-01-02")

	rows := sqlmock.NewRows([]string{"id", "user_id", "service", "date", "time", "client_name", "phone", "status"}).
		AddRow(1, 123456789, "Маникюр", today, "10:00", "Test Client", "+79991234567", "pending").
		AddRow(2, 987654321, "Педикюр", threeDaysLater, "11:00", "Another Client", "+79991234568", "pending")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "bookings" WHERE date >= $1 AND date <= $2 AND status != $3 ORDER BY date ASC, time ASC`)).
		WithArgs(today, threeDaysLater, "cancelled").
		WillReturnRows(rows)

	reportService := &ReportService{}
	bookings, err := reportService.GetBookingsForNext3Days()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(bookings) != 2 {
		t.Errorf("Expected 2 bookings, got %d", len(bookings))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

// Тест для GetBookingsForNext3Days — нет записей
func TestGetBookingsForNext3Days_NoBookings(t *testing.T) {
	_, mock := setupTestDBForReport(t)

	today := time.Now().Format("2006-01-02")
	threeDaysLater := time.Now().AddDate(0, 0, 3).Format("2006-01-02")

	rows := sqlmock.NewRows([]string{"id", "user_id", "service", "date", "time", "client_name", "phone", "status"})

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "bookings" WHERE date >= $1 AND date <= $2 AND status != $3 ORDER BY date ASC, time ASC`)).
		WithArgs(today, threeDaysLater, "cancelled").
		WillReturnRows(rows)

	reportService := &ReportService{}
	bookings, err := reportService.GetBookingsForNext3Days()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(bookings) != 0 {
		t.Errorf("Expected 0 bookings, got %d", len(bookings))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

// Тест для GenerateReport — с записями
func TestGenerateReport_WithBookings(t *testing.T) {
	bookings := []models.Booking{
		{
			ID:         1,
			UserID:     123456789,
			Service:    "Маникюр",
			Date:       time.Now().Format("2006-01-02"),
			Time:       "10:00",
			ClientName: "Test Client",
			Phone:      "+79991234567",
			Status:     "pending",
		},
		{
			ID:         2,
			UserID:     987654321,
			Service:    "Педикюр",
			Date:       time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
			Time:       "11:00",
			ClientName: "Another Client",
			Phone:      "+79991234568",
			Status:     "pending",
		},
	}

	reportService := &ReportService{}
	data, err := reportService.GenerateReport(bookings)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(data) == 0 {
		t.Error("Expected PDF data, got empty")
	}
	if len(data) < 100 {
		t.Errorf("Expected PDF data size > 100 bytes, got %d", len(data))
	}
}

// Тест для GenerateReport — без записей
func TestGenerateReport_NoBookings(t *testing.T) {
	bookings := []models.Booking{}

	reportService := &ReportService{}
	data, err := reportService.GenerateReport(bookings)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(data) == 0 {
		t.Error("Expected PDF data, got empty")
	}
	if len(data) < 100 {
		t.Errorf("Expected PDF data size > 100 bytes, got %d", len(data))
	}
}
