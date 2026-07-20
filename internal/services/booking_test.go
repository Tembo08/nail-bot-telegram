package services

import (
	"regexp"
	"testing"

	"nail_bot/internal/storage"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupTestDB создаёт мок-базу данных для тестов
func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

	// Сохраняем в глобальную переменную
	storage.DB = gormDB

	return gormDB, mock
}

// Тест для CreateBooking
func TestCreateBooking(t *testing.T) {
	tests := []struct {
		name          string
		userID        int64
		service       string
		date          string
		timeStr       string
		clientName    string
		phone         string
		exists        bool
		expectError   bool
		errorContains string
	}{
		{
			name:        "Successful booking",
			userID:      123456789,
			service:     "Маникюр",
			date:        "2026-07-22",
			timeStr:     "10:00",
			clientName:  "Test Client",
			phone:       "+79991234567",
			exists:      false,
			expectError: false,
		},
		{
			name:          "Time already booked",
			userID:        123456789,
			service:       "Маникюр",
			date:          "2026-07-22",
			timeStr:       "10:00",
			clientName:    "Test Client",
			phone:         "+79991234567",
			exists:        true,
			expectError:   true,
			errorContains: "это время уже занято",
		},
	}

	bookingService := &BookingService{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, mock := setupTestDB(t)

			// Мокаем CheckAvailability
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "bookings" WHERE date = $1 AND time = $2 AND status != $3`)).
				WithArgs(tt.date, tt.timeStr, "cancelled").
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(func() int {
					if tt.exists {
						return 1
					}
					return 0
				}()))

			if !tt.exists {
				// Мокаем Create
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "bookings"`)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			}

			booking, err := bookingService.CreateBooking(
				tt.userID,
				tt.service,
				tt.date,
				tt.timeStr,
				tt.clientName,
				tt.phone,
			)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				if tt.errorContains != "" && err != nil {
					if !contains(err.Error(), tt.errorContains) {
						t.Errorf("Expected error containing '%s', got '%s'", tt.errorContains, err.Error())
					}
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if booking == nil {
					t.Error("Expected booking, got nil")
				}
			}

			// Проверяем, что все ожидания выполнены
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

// Тест для CheckAvailability
func TestCheckAvailability(t *testing.T) {
	tests := []struct {
		name      string
		date      string
		timeStr   string
		count     int
		expectErr bool
		want      bool
	}{
		{
			name:    "Time is available",
			date:    "2026-07-22",
			timeStr: "10:00",
			count:   0,
			want:    false,
		},
		{
			name:    "Time is booked",
			date:    "2026-07-22",
			timeStr: "10:00",
			count:   1,
			want:    true,
		},
	}

	bookingService := &BookingService{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, mock := setupTestDB(t)

			mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "bookings" WHERE date = $1 AND time = $2 AND status != $3`)).
				WithArgs(tt.date, tt.timeStr, "cancelled").
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(tt.count))

			got, err := bookingService.CheckAvailability(tt.date, tt.timeStr)

			if tt.expectErr && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if got != tt.want {
				t.Errorf("CheckAvailability() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

// Тест для GetUserBookings
func TestGetUserBookings(t *testing.T) {
	_, mock := setupTestDB(t)

	userID := int64(123456789)

	rows := sqlmock.NewRows([]string{"id", "user_id", "service", "date", "time", "client_name", "phone", "status"}).
		AddRow(1, userID, "Маникюр", "2026-07-22", "10:00", "Test Client", "+79991234567", "pending")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "bookings" WHERE user_id = $1 AND status != $2 ORDER BY date ASC, time ASC`)).
		WithArgs(userID, "cancelled").
		WillReturnRows(rows)

	bookingService := &BookingService{}
	bookings, err := bookingService.GetUserBookings(userID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(bookings) != 1 {
		t.Errorf("Expected 1 booking, got %d", len(bookings))
	}
	if bookings[0].Service != "Маникюр" {
		t.Errorf("Expected 'Маникюр', got '%s'", bookings[0].Service)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

// Тест для GetAllBookings
func TestGetAllBookings(t *testing.T) {
	_, mock := setupTestDB(t)

	rows := sqlmock.NewRows([]string{"id", "user_id", "service", "date", "time", "client_name", "phone", "status"}).
		AddRow(1, 123456789, "Маникюр", "2026-07-22", "10:00", "Test Client", "+79991234567", "pending").
		AddRow(2, 987654321, "Педикюр", "2026-07-23", "11:00", "Another Client", "+79991234568", "pending")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "bookings" WHERE status != $1 ORDER BY date ASC, time ASC`)).
		WithArgs("cancelled").
		WillReturnRows(rows)

	bookingService := &BookingService{}
	bookings, err := bookingService.GetAllBookings()

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

// Тест для CancelBooking
func TestCancelBooking(t *testing.T) {
	tests := []struct {
		name         string
		bookingID    uint
		userID       int64
		rowsAffected int64
		expectError  bool
	}{
		{
			name:         "Successfully cancel booking",
			bookingID:    1,
			userID:       123456789,
			rowsAffected: 1,
			expectError:  false,
		},
		{
			name:         "Booking not found",
			bookingID:    99,
			userID:       123456789,
			rowsAffected: 0,
			expectError:  true,
		},
	}

	bookingService := &BookingService{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, mock := setupTestDB(t)

			mock.ExpectBegin()
			// Исправляем regexp для соответствия реальному запросу с updated_at
			mock.ExpectExec(regexp.QuoteMeta(`UPDATE "bookings" SET "status"=$1,"updated_at"=$2 WHERE id = $3 AND user_id = $4`)).
				WithArgs("cancelled", sqlmock.AnyArg(), tt.bookingID, tt.userID).
				WillReturnResult(sqlmock.NewResult(0, tt.rowsAffected))
			mock.ExpectCommit()

			err := bookingService.CancelBooking(tt.bookingID, tt.userID)

			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

// Тест для GetActiveBookingsByUserID
func TestGetActiveBookingsByUserID(t *testing.T) {
	_, mock := setupTestDB(t)

	userID := int64(123456789)

	rows := sqlmock.NewRows([]string{"id", "user_id", "service", "date", "time", "client_name", "phone", "status"}).
		AddRow(1, userID, "Маникюр", "2026-07-22", "10:00", "Test Client", "+79991234567", "pending")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "bookings" WHERE user_id = $1 AND status NOT IN ($2,$3) ORDER BY date ASC, time ASC`)).
		WithArgs(userID, "cancelled", "completed").
		WillReturnRows(rows)

	bookingService := &BookingService{}
	bookings, err := bookingService.GetActiveBookingsByUserID(userID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(bookings) != 1 {
		t.Errorf("Expected 1 booking, got %d", len(bookings))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

// Тест для CancelBookingByID
func TestCancelBookingByID(t *testing.T) {
	tests := []struct {
		name         string
		bookingID    uint
		userID       int64
		rowsAffected int64
		expectError  bool
	}{
		{
			name:         "Successfully cancel booking by ID",
			bookingID:    1,
			userID:       123456789,
			rowsAffected: 1,
			expectError:  false,
		},
		{
			name:         "Booking not found by ID",
			bookingID:    99,
			userID:       123456789,
			rowsAffected: 0,
			expectError:  true,
		},
	}

	bookingService := &BookingService{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, mock := setupTestDB(t)

			mock.ExpectBegin()
			// Исправляем regexp для соответствия реальному запросу с updated_at
			mock.ExpectExec(regexp.QuoteMeta(`UPDATE "bookings" SET "status"=$1,"updated_at"=$2 WHERE id = $3 AND user_id = $4`)).
				WithArgs("cancelled", sqlmock.AnyArg(), tt.bookingID, tt.userID).
				WillReturnResult(sqlmock.NewResult(0, tt.rowsAffected))
			mock.ExpectCommit()

			err := bookingService.CancelBookingByID(tt.bookingID, tt.userID)

			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

// Вспомогательная функция для проверки наличия подстроки
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
