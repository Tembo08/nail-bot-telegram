package handlers

import (
	"nail_bot/internal/keyboards"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageHandler struct {
	Bot            *tgbotapi.BotAPI
	BookingHandler *BookingHandler // ← добавляем это поле
}

func (h *MessageHandler) Start(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "👋 Добро пожаловать в Nail Bot!\n\n"+
		"Я помогу вам записаться на маникюрные услуги.\n\n"+
		"Выберите действие из меню:")

	msg.ReplyMarkup = keyboards.MainMenu()
	_, err := h.Bot.Send(msg)
	return err
}

func (h *MessageHandler) HandleMessage(message *tgbotapi.Message) error {
	chatID := message.Chat.ID
	text := message.Text

	// Проверяем, есть ли активная сессия в состоянии contact
	if h.BookingHandler != nil {
		session, _ := h.BookingHandler.SessionService.GetSession(chatID)
		if session != nil && session.Step == "contact" {
			// Проверяем, не является ли сообщение командой или кнопкой
			if strings.HasPrefix(text, "/") ||
				text == "📅 Записаться" ||
				text == "📋 Мои записи" ||
				text == "ℹ️ О нас" ||
				text == "📞 Контакты" {
				// Если это команда — отменяем сессию и обрабатываем как обычно
				h.BookingHandler.SessionService.DeleteSession(chatID)
				// Продолжаем выполнение (выходим из условия)
			} else {
				// Это ввод контакта
				return h.BookingHandler.HandleContactInput(chatID, text)
			}
		}
	}

	switch text {
	case "📅 Записаться":
		return h.BookingHandler.StartBooking(chatID)
	case "📋 Мои записи":
		// Отменяем сессию, если она есть
		h.BookingHandler.SessionService.DeleteSession(chatID)
		return h.BookingHandler.ShowMyBookings(chatID)
	case "ℹ️ О нас":
		h.BookingHandler.SessionService.DeleteSession(chatID)
		msg := tgbotapi.NewMessage(chatID, "💅 Nail Studio — профессиональный маникюрный салон.\n\n"+
			"🕐 9:00-21:00\n"+
			"📍 ул. Красивая, д. 1\n"+
			"👤 Мастер: Инесса")
		_, err := h.Bot.Send(msg)
		return err
	case "📞 Контакты":
		h.BookingHandler.SessionService.DeleteSession(chatID)
		msg := tgbotapi.NewMessage(chatID, "📞 Контакты:\n\n"+
			"📱 +7 (999) 123-45-67\n"+
			"📍 ул. Красивая, д. 1\n"+
			"🌐 @nail_studio")
		_, err := h.Bot.Send(msg)
		return err
	default:
		// Если пользователь ввёл что-то другое — показываем меню
		msg := tgbotapi.NewMessage(chatID, "🔽 Используйте кнопки меню:")
		msg.ReplyMarkup = keyboards.MainMenu()
		_, err := h.Bot.Send(msg)
		return err
	}
}
