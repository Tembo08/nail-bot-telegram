package keyboards

import (
	"fmt"
	"nail_bot/internal/models"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Главное меню
func MainMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📅 Записаться"),
			tgbotapi.NewKeyboardButton("📋 Мои записи"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("ℹ️ О нас"),
			tgbotapi.NewKeyboardButton("📞 Контакты"),
		),
	)
}

// Клавиатура с услугами
func ServicesMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💅 Маникюр", "service_manicure"),
			tgbotapi.NewInlineKeyboardButtonData("💅 Педикюр", "service_pedicure"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💅 Дизайн", "service_design"),
			tgbotapi.NewInlineKeyboardButtonData("💅 Наращивание", "service_extension"),
		),
	)
}

// Клавиатура с датами (14 дней, без воскресений)
func DateKeyboard() tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	now := time.Now()

	for i := 0; i < 14; i++ {
		date := now.AddDate(0, 0, i)
		if date.Weekday() == time.Sunday {
			continue
		}
		label := date.Format("02.01 (Mon)")
		data := date.Format("2006-01-02")
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(label, "date_"+data),
		))
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "back_service"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// Клавиатура с временем (10:00 - 20:00, шаг 30 мин)
func TimeKeyboard() tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	for hour := 10; hour <= 20; hour++ {
		for min := 0; min < 60; min += 30 {
			timeStr := time.Date(0, 0, 0, hour, min, 0, 0, time.UTC).Format("15:04")
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(timeStr, "time_"+timeStr))

			if len(row) == 3 {
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(row...))
				row = []tgbotapi.InlineKeyboardButton{}
			}
		}
	}

	if len(row) > 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(row...))
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "back_date"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// Клавиатура подтверждения
func ConfirmKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Подтвердить", "confirm_yes"),
			tgbotapi.NewInlineKeyboardButtonData("❌ Отменить", "confirm_no"),
		),
	)
}

// 12 июля
// BookingKeyboard создаёт inline-клавиатуру с записями
func BookingKeyboard(bookings []models.Booking) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	if len(bookings) == 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📭 Нет активных записей", "no_bookings"),
		))
	}

	for _, b := range bookings {
		label := fmt.Sprintf("%s | %s %s", b.Service, b.Date, b.Time)
		data := fmt.Sprintf("cancel_%d", b.ID)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(label, data),
		))
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "back_to_main"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// CancelConfirmKeyboard клавиатура подтверждения отмены
func CancelConfirmKeyboard(bookingID uint) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Да, отменить", fmt.Sprintf("confirm_cancel_%d", bookingID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Нет, не надо", "cancel_no"),
		),
	)
}
