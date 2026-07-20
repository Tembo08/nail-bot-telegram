package services

import (
	"regexp"
	"testing"
	"time"

	"nail_bot/internal/storage"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupTestDB для reminder_test
func setupTestDBForReminder(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

// Тест для GetBookingsForTomorrow
func TestGetBookingsForTomorrow(t *testing.T) {
	tests := []struct {
		name      string
		mockRows  *sqlmock.Rows
		expectErr bool
		wantCount int
	}{
		{
			name: "Bookings found for tomorrow",
			mockRows: sqlmock.NewRows([]string{"id", "user_id", "service", "date", "time", "client_name", "phone", "status"}).
				AddRow(1, 123456789, "Маникюр", time.Now().AddDate(0, 0, 1).Format("2006-01-02"), "10:00", "Test Client", "+79991234567", "pending").
				AddRow(2, 987654321, "Педикюр", time.Now().AddDate(0, 0, 1).Format("2006-01-02"), "11:00", "Another Client", "+79991234568", "pending"),
			expectErr: false,
			wantCount: 2,
		},
		{
			name:      "No bookings for tomorrow",
			mockRows:  sqlmock.NewRows([]string{"id", "user_id", "service", "date", "time", "client_name", "phone", "status"}),
			expectErr: false,
			wantCount: 0,
		},
	}

	reminderService := &ReminderService{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, mock := setupTestDBForReminder(t)

			tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

			mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "bookings" WHERE date = $1 AND status != $2`)).
				WithArgs(tomorrow, "cancelled").
				WillReturnRows(tt.mockRows)

			bookings, err := reminderService.GetBookingsForTomorrow()

			if tt.expectErr && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if len(bookings) != tt.wantCount {
				t.Errorf("Expected %d bookings, got %d", tt.wantCount, len(bookings))
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

// Тест для формата даты завтра
func TestTomorrowDate(t *testing.T) {
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	if tomorrow == "" {
		t.Error("Tomorrow date is empty")
	}

	// Проверяем формат YYYY-MM-DD
	if len(tomorrow) != 10 {
		t.Errorf("Tomorrow date format is incorrect: %s", tomorrow)
	}
}

// Тест для проверки, что сегодня != завтра
func TestTomorrowIsNotToday(t *testing.T) {
	today := time.Now().Format("2006-01-02")
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	if today == tomorrow {
		t.Error("Today and tomorrow are the same")
	}
}
