package handlers

import (
	"fmt"
	"log"
	"nail_bot/internal/keyboards"
	"nail_bot/internal/models"
	"nail_bot/internal/services"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BookingHandler struct {
	Bot            *tgbotapi.BotAPI
	BookingService *services.BookingService
	SessionService *services.SessionService
}

func NewBookingHandler(bot *tgbotapi.BotAPI) *BookingHandler {
	return &BookingHandler{
		Bot:            bot,
		BookingService: &services.BookingService{},
		SessionService: &services.SessionService{},
	}
}

// Начало записи
func (h *BookingHandler) StartBooking(chatID int64) error {
	// Создаём или получаем сессию
	_, err := h.SessionService.GetOrCreateSession(chatID)
	if err != nil {
		log.Printf("Ошибка создания сессии: %v", err)
		return err
	}

	// Обновляем шаг
	h.SessionService.UpdateSession(chatID, models.StepService, nil)

	msg := tgbotapi.NewMessage(chatID, "💅 Выберите услугу:")
	msg.ReplyMarkup = keyboards.ServicesMenu()
	_, err = h.Bot.Send(msg)
	return err
}

// Обработка callback-запросов
func (h *BookingHandler) HandleBookingCallback(callback *tgbotapi.CallbackQuery) error {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	h.Bot.Send(tgbotapi.NewCallback(callback.ID, ""))

	switch {
	case strings.HasPrefix(data, "service_"):
		return h.handleService(chatID, data)
	case strings.HasPrefix(data, "date_"):
		return h.handleDate(chatID, data)
	case strings.HasPrefix(data, "time_"):
		return h.handleTime(chatID, data)
	case data == "back_service":
		return h.goToService(chatID)
	case data == "back_date":
		return h.goToDate(chatID)
	case data == "confirm_yes":
		return h.confirmBooking(chatID)
	case data == "confirm_no":
		return h.cancelBooking(chatID)

		// 12 июля
		// В функции HandleBookingCallback добавьте:
	case strings.HasPrefix(data, "cancel_"):
		return h.HandleCancelBooking(chatID, data)
	case strings.HasPrefix(data, "confirm_cancel_"):
		return h.ConfirmCancelBooking(chatID, data)
	case data == "cancel_no":
		msg := tgbotapi.NewMessage(chatID, "✅ Отмена отменена!")
		msg.ReplyMarkup = keyboards.MainMenu()
		_, err := h.Bot.Send(msg)
		return err
	case data == "back_to_main":
		msg := tgbotapi.NewMessage(chatID, "🔙 Возвращаемся в главное меню")
		msg.ReplyMarkup = keyboards.MainMenu()
		_, err := h.Bot.Send(msg)
		return err
	default:
		return h.handleUnknown(chatID)
	}
}

func (h *BookingHandler) handleService(chatID int64, data string) error {
	serviceMap := map[string]string{
		"manicure":  "Маникюр",
		"pedicure":  "Педикюр",
		"design":    "Дизайн",
		"extension": "Наращивание",
	}

	service, ok := serviceMap[strings.TrimPrefix(data, "service_")]
	if !ok {
		msg := tgbotapi.NewMessage(chatID, "❌ Неизвестная услуга")
		_, err := h.Bot.Send(msg)
		return err
	}

	h.SessionService.UpdateSession(chatID, models.StepDate, map[string]string{"service": service})

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(
		"✅ Вы выбрали: *%s*\n\n📅 Выберите дату (пн-сб):",
		service,
	))
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboards.DateKeyboard()
	_, err := h.Bot.Send(msg)
	return err
}

func (h *BookingHandler) handleDate(chatID int64, data string) error {
	date := strings.TrimPrefix(data, "date_")
	h.SessionService.UpdateSession(chatID, models.StepTime, map[string]string{"date": date})

	session, _ := h.SessionService.GetSession(chatID)

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(
		"✅ Вы выбрали дату: *%s*\n\n⏰ Выберите время (10:00-20:00):",
		session.Date,
	))
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboards.TimeKeyboard()
	_, err := h.Bot.Send(msg)
	return err
}

func (h *BookingHandler) handleTime(chatID int64, data string) error {
	timeStr := strings.TrimPrefix(data, "time_")
	session, _ := h.SessionService.GetSession(chatID)

	// Проверяем занятость
	available, err := h.BookingService.CheckAvailability(session.Date, timeStr)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "⚠️ Ошибка проверки времени")
		_, err := h.Bot.Send(msg)
		return err
	}

	if available {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(
			"❌ Время *%s* уже занято. Выберите другое:",
			timeStr,
		))
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = keyboards.TimeKeyboard()
		_, err := h.Bot.Send(msg)
		return err
	}

	h.SessionService.UpdateSession(chatID, models.StepContact, map[string]string{"time": timeStr})

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(
		"✅ Вы выбрали время: *%s*\n\n📝 Введите *имя* и *телефон*:\nПример: *Иван, +7 (999) 123-45-67*",
		timeStr,
	))
	msg.ParseMode = "Markdown"
	_, err = h.Bot.Send(msg)
	return err
}

func (h *BookingHandler) HandleContactInput(chatID int64, text string) error {
	session, err := h.SessionService.GetSession(chatID)
	if err != nil || session.Step != models.StepContact {
		return nil
	}

	parts := strings.SplitN(text, ",", 2)
	if len(parts) != 2 {
		msg := tgbotapi.NewMessage(chatID, "❌ Неправильный формат.\n\nВведите: *Имя, Телефон*")
		msg.ParseMode = "Markdown"
		_, err := h.Bot.Send(msg)
		return err
	}

	h.SessionService.UpdateSession(chatID, models.StepConfirm, map[string]string{"contact": text})

	session, _ = h.SessionService.GetSession(chatID)

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(
		"📋 *Проверьте данные:*\n\n"+
			"💅 Услуга: %s\n"+
			"📅 Дата: %s\n"+
			"⏰ Время: %s\n"+
			"👤 Мастер: %s\n"+
			"📝 Контакт: %s\n\n"+
			"Всё верно?",
		session.Service, session.Date, session.Time, services.MasterName, session.Contact,
	))
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboards.ConfirmKeyboard()
	_, err = h.Bot.Send(msg)
	return err
}

func (h *BookingHandler) confirmBooking(chatID int64) error {
	session, err := h.SessionService.GetSession(chatID)
	if err != nil {
		return err
	}

	parts := strings.SplitN(session.Contact, ",", 2)
	name := strings.TrimSpace(parts[0])
	phone := strings.TrimSpace(parts[1])

	booking, err := h.BookingService.CreateBooking(
		chatID, session.Service, session.Date, session.Time, name, phone,
	)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ %s", err.Error()))
		_, err := h.Bot.Send(msg)
		return err
	}

	h.SessionService.DeleteSession(chatID)

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(
		"✅ *Запись подтверждена!*\n\n"+
			"💅 Услуга: %s\n"+
			"📅 Дата: %s\n"+
			"⏰ Время: %s\n"+
			"👤 Мастер: %s\n\n"+
			"📍 ул. Красивая, д. 1\n"+
			"📱 +7 (999) 123-45-67",
		session.Service, session.Date, session.Time, services.MasterName,
	))
	msg.ParseMode = "Markdown"
	_, err = h.Bot.Send(msg)

	log.Printf("📝 Новая запись #%d: %s, %s, %s", booking.ID, name, session.Service, session.Date)
	return err
}

func (h *BookingHandler) cancelBooking(chatID int64) error {
	h.SessionService.DeleteSession(chatID)

	msg := tgbotapi.NewMessage(chatID, "❌ Запись отменена.")
	msg.ReplyMarkup = keyboards.MainMenu()
	_, err := h.Bot.Send(msg)
	return err
}

func (h *BookingHandler) goToService(chatID int64) error {
	h.SessionService.UpdateSession(chatID, models.StepService, nil)

	msg := tgbotapi.NewMessage(chatID, "💅 Выберите услугу:")
	msg.ReplyMarkup = keyboards.ServicesMenu()
	_, err := h.Bot.Send(msg)
	return err
}

func (h *BookingHandler) goToDate(chatID int64) error {
	h.SessionService.UpdateSession(chatID, models.StepDate, nil)

	msg := tgbotapi.NewMessage(chatID, "📅 Выберите дату:")
	msg.ReplyMarkup = keyboards.DateKeyboard()
	_, err := h.Bot.Send(msg)
	return err
}

func (h *BookingHandler) handleUnknown(chatID int64) error {
	msg := tgbotapi.NewMessage(chatID, "❌ Неизвестная команда")
	_, err := h.Bot.Send(msg)
	return err
}

// 12 июля
// ShowMyBookings показывает записи пользователя
func (h *BookingHandler) ShowMyBookings(chatID int64) error {
	bookings, err := h.BookingService.GetActiveBookingsByUserID(chatID)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "❌ Ошибка при получении записей")
		_, err := h.Bot.Send(msg)
		return err
	}

	if len(bookings) == 0 {
		msg := tgbotapi.NewMessage(chatID, "📭 У вас пока нет активных записей.")
		msg.ReplyMarkup = keyboards.MainMenu()
		_, err := h.Bot.Send(msg)
		return err
	}

	msg := tgbotapi.NewMessage(chatID, "📋 *Ваши записи:*\n\nНажмите на запись, чтобы отменить её.")
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboards.BookingKeyboard(bookings)
	_, err = h.Bot.Send(msg)
	return err
}

// HandleCancelBooking обрабатывает отмену записи
func (h *BookingHandler) HandleCancelBooking(chatID int64, data string) error {
	// Парсим ID записи
	var bookingID uint
	fmt.Sscanf(strings.TrimPrefix(data, "cancel_"), "%d", &bookingID)

	if bookingID == 0 {
		msg := tgbotapi.NewMessage(chatID, "❌ Неверный ID записи")
		_, err := h.Bot.Send(msg)
		return err
	}

	// Показываем подтверждение
	msg := tgbotapi.NewMessage(chatID, "⚠️ *Вы уверены, что хотите отменить эту запись?*\n\nЭто действие нельзя отменить.")
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboards.CancelConfirmKeyboard(bookingID)
	_, err := h.Bot.Send(msg)
	return err
}

// ConfirmCancelBooking подтверждает отмену записи
func (h *BookingHandler) ConfirmCancelBooking(chatID int64, data string) error {
	var bookingID uint
	fmt.Sscanf(strings.TrimPrefix(data, "confirm_cancel_"), "%d", &bookingID)

	if bookingID == 0 {
		msg := tgbotapi.NewMessage(chatID, "❌ Неверный ID записи")
		_, err := h.Bot.Send(msg)
		return err
	}

	err := h.BookingService.CancelBookingByID(bookingID, chatID)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ %s", err.Error()))
		_, err := h.Bot.Send(msg)
		return err
	}

	msg := tgbotapi.NewMessage(chatID, "✅ *Запись успешно отменена!*")
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboards.MainMenu()
	_, err = h.Bot.Send(msg)
	return err
}

// 13 июля
// SendReport отправляет PDF-отчёт админу
func (h *BookingHandler) SendReport(chatID int64) error {
	reportService := &services.ReportService{}

	// Получаем записи
	bookings, err := reportService.GetBookingsForNext3Days()
	if err != nil {
		log.Printf("❌ Ошибка при получении записей: %v", err)
		msg := tgbotapi.NewMessage(chatID, "❌ Ошибка при получении записей")
		_, err := h.Bot.Send(msg)
		return err
	}

	log.Printf("📊 Найдено записей для отчёта: %d", len(bookings))

	// Генерируем PDF
	pdfData, err := reportService.GenerateReport(bookings)
	if err != nil {
		log.Printf("❌ Ошибка при генерации отчёта: %v", err)
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ Ошибка при генерации отчёта: %v", err))
		_, err := h.Bot.Send(msg)
		return err
	}

	log.Printf("✅ PDF сгенерирован, размер: %d байт", len(pdfData))

	// Отправляем PDF
	file := tgbotapi.FileBytes{
		Name:  "report.pdf",
		Bytes: pdfData,
	}

	msg := tgbotapi.NewDocument(chatID, file)
	msg.Caption = fmt.Sprintf("📊 Отчёт на 3 дня\nВсего записей: %d", len(bookings))

	_, err = h.Bot.Send(msg)
	if err != nil {
		log.Printf("❌ Ошибка при отправке PDF: %v", err)
		return err
	}

	log.Printf("✅ Отчёт отправлен успешно")
	return nil
}

// 13.07
// ShowAllBookings показывает все записи админу
func (h *BookingHandler) ShowAllBookings(chatID int64) error {
	bookings, err := h.BookingService.GetAllBookings()
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "❌ Ошибка при получении записей")
		_, err := h.Bot.Send(msg)
		return err
	}

	if len(bookings) == 0 {
		msg := tgbotapi.NewMessage(chatID, "📭 Нет активных записей.")
		_, err := h.Bot.Send(msg)
		return err
	}

	var text string
	text = "📋 *Все записи:*\n\n"
	for i, b := range bookings {
		text += fmt.Sprintf("%d. *%s* | %s %s\n   👤 %s | 📱 %s\n   🆔 %d\n\n",
			i+1, b.Service, b.Date, b.Time, b.ClientName, b.Phone, b.ID)
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	_, err = h.Bot.Send(msg)
	return err
}
