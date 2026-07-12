package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"nail_bot/configs"
	"nail_bot/internal/handlers"
	"nail_bot/internal/storage"

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
