package utils

import (
	"regexp"
	"strings"
)

// IsValidPhone проверяет, что телефон соответствует формату
// Поддерживает: +7 (999) 123-45-67, +79991234567, 89991234567
// IsValidPhone проверяет, что телефон соответствует формату
func IsValidPhone(phone string) bool {
	// Проверяем, что строка не пустая
	if len(phone) == 0 {
		return false
	}

	// Проверяем, что все символы допустимы (цифры, +, пробелы, скобки, дефисы)
	for _, ch := range phone {
		isDigit := ch >= '0' && ch <= '9'
		isPlus := ch == '+'
		isSpace := ch == ' '
		isBracket := ch == '(' || ch == ')'
		isDash := ch == '-'
		if !isDigit && !isPlus && !isSpace && !isBracket && !isDash {
			return false
		}
	}

	// Убираем все лишние символы, оставляем только цифры и +
	cleaned := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' || r == '+' {
			return r
		}
		return -1
	}, phone)

	// Проверяем длину и начало
	if len(cleaned) < 10 || len(cleaned) > 12 {
		return false
	}

	if !strings.HasPrefix(cleaned, "7") &&
		!strings.HasPrefix(cleaned, "8") &&
		!strings.HasPrefix(cleaned, "+7") {
		return false
	}

	return true
}

// IsValidTime проверяет, что время в формате HH:MM
func IsValidTime(timeStr string) bool {
	re := regexp.MustCompile(`^([01][0-9]|2[0-3]):([0-5][0-9])$`)
	return re.MatchString(timeStr)
}

// IsValidDate проверяет, что дата в формате YYYY-MM-DD
func IsValidDate(dateStr string) bool {
	re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	return re.MatchString(dateStr)
}
