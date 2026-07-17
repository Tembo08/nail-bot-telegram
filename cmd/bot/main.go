package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"nail_bot/configs"
	"nail_bot/internal/handlers"
	"nail_bot/internal/services"
	"nail_bot/internal/storage"

	"github.com/robfig/cron/v3"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg := configs.LoadConfig()

	if cfg.TelegramBotToken == "" {
		log.Fatal("❌ TELEGRAM_BOT_TOKEN не найден!")
	}

	if err := storage.InitDB(cfg); err != nil {
		log.Fatal("❌ Ошибка подключения к БД:", err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		log.Fatal("❌ Ошибка создания бота:", err)
	}

	bot.Debug = true
	log.Printf("✅ Бот запущен: @%s", bot.Self.UserName)

	// Запускаем планировщик
	startScheduler(bot, cfg.AdminID)

	bookingHandler := handlers.NewBookingHandler(bot)

	messageHandler := &handlers.MessageHandler{
		Bot:            bot,
		BookingHandler: bookingHandler,
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("🛑 Бот остановлен")
		os.Exit(0)
	}()

	log.Println("⏳ Ожидание сообщений...")

	for update := range updates {
		if update.Message != nil && update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				messageHandler.Start(update.Message)
			default:
				// Все остальные команды (включая /admin) отправляем в HandleMessage
				messageHandler.HandleMessage(update.Message)
			}
			continue
		}

		if update.Message != nil {
			messageHandler.HandleMessage(update.Message)
			continue
		}

		if update.CallbackQuery != nil {
			bookingHandler.HandleBookingCallback(update.CallbackQuery)
		}
	}
}

// 13.07
func startScheduler(bot *tgbotapi.BotAPI, adminID int64) {
	c := cron.New()

	// Каждый день в 9:00 отправляем напоминания о записях на завтра
	c.AddFunc("0 9 * * *", func() {
		reminderService := &services.ReminderService{}
		bookings, err := reminderService.GetBookingsForTomorrow()
		if err != nil {
			log.Printf("Ошибка получения записей для напоминаний: %v", err)
			return
		}

		for _, booking := range bookings {
			msg := tgbotapi.NewMessage(
				booking.UserID,
				fmt.Sprintf("🔔 *Напоминание о записи!*\n\n"+
					"💅 Услуга: %s\n"+
					"📅 Дата: %s\n"+
					"⏰ Время: %s\n"+
					"👤 Мастер: %s\n\n"+
					"📍 ул. Красивая, д. 1\n"+
					"Ждём вас! 💅",
					booking.Service, booking.Date, booking.Time, services.MasterName,
				),
			)
			msg.ParseMode = "Markdown"
			bot.Send(msg)
		}

		log.Printf("✅ Отправлено %d напоминаний", len(bookings))
	})

	c.Start()
	log.Println("⏰ Планировщик напоминаний запущен")
}
