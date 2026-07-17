package utils

import "testing"

func TestIsValidPhone(t *testing.T) {
	tests := []struct {
		name  string
		phone string
		want  bool
	}{
		// Валидные номера
		{"Valid with +7 and spaces", "+7 (999) 123-45-67", true},
		{"Valid with +7 without spaces", "+79991234567", true},
		{"Valid with 8", "89991234567", true},
		{"Valid with 7", "79991234567", true},
		{"Valid short", "7999123456", true},
		{"Valid long", "799912345678", true},

		// Невалидные номера
		{"Invalid too short", "123", false},
		{"Invalid letters", "abc", false},
		{"Invalid wrong prefix", "9991234567", false},
		{"Invalid empty", "", false},
		{"Invalid with letters", "+7 (999) 123-45-6A", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidPhone(tt.phone); got != tt.want {
				t.Errorf("IsValidPhone(%q) = %v, want %v", tt.phone, got, tt.want)
			}
		})
	}
}

func TestIsValidTime(t *testing.T) {
	tests := []struct {
		name string
		time string
		want bool
	}{
		{"Valid 10:00", "10:00", true},
		{"Valid 23:59", "23:59", true},
		{"Valid 00:00", "00:00", true},
		{"Invalid 24:00", "24:00", false},
		{"Invalid 10:60", "10:60", false},
		{"Invalid abc", "abc", false},
		{"Invalid 10:5", "10:5", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidTime(tt.time); got != tt.want {
				t.Errorf("IsValidTime(%q) = %v, want %v", tt.time, got, tt.want)
			}
		})
	}
}

func TestIsValidDate(t *testing.T) {
	tests := []struct {
		name string
		date string
		want bool
	}{
		{"Valid date", "2026-07-15", true},
		{"Invalid date", "2026/07/15", false},
		{"Invalid format", "15-07-2026", false},
		{"Invalid short", "2026-7-15", false},
		{"Invalid empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidDate(tt.date); got != tt.want {
				t.Errorf("IsValidDate(%q) = %v, want %v", tt.date, got, tt.want)
			}
		})
	}
}

