package models

import (
	"testing"
	"time"
)

func TestUserTableName(t *testing.T) {
	user := User{}
	expected := "users"
	if got := user.TableName(); got != expected {
		t.Errorf("User.TableName() = %v, want %v", got, expected)
	}
}

func TestBookingTableName(t *testing.T) {
	booking := Booking{}
	expected := "bookings"
	if got := booking.TableName(); got != expected {
		t.Errorf("Booking.TableName() = %v, want %v", got, expected)
	}
}

func TestUserSessionTableName(t *testing.T) {
	session := UserSession{}
	expected := "user_sessions"
	if got := session.TableName(); got != expected {
		t.Errorf("UserSession.TableName() = %v, want %v", got, expected)
	}
}

func TestUserFields(t *testing.T) {
	user := User{
		ID:         1,
		TelegramID: 123456789,
		FirstName:  "Test",
		LastName:   "User",
		Username:   "testuser",
		Phone:      "+79991234567",
		CreatedAt:  time.Now(),
	}

	if user.ID != 1 {
		t.Errorf("User.ID = %v, want 1", user.ID)
	}
	if user.TelegramID != 123456789 {
		t.Errorf("User.TelegramID = %v, want 123456789", user.TelegramID)
	}
	if user.FirstName != "Test" {
		t.Errorf("User.FirstName = %v, want 'Test'", user.FirstName)
	}
}

func TestBookingFields(t *testing.T) {
	booking := Booking{
		ID:         1,
		UserID:     123456789,
		Service:    "Маникюр",
		Date:       "2026-07-15",
		Time:       "10:00",
		ClientName: "Test Client",
		Phone:      "+79991234567",
		Status:     "pending",
	}

	if booking.ID != 1 {
		t.Errorf("Booking.ID = %v, want 1", booking.ID)
	}
	if booking.Service != "Маникюр" {
		t.Errorf("Booking.Service = %v, want 'Маникюр'", booking.Service)
	}
	if booking.Status != "pending" {
		t.Errorf("Booking.Status = %v, want 'pending'", booking.Status)
	}
}
