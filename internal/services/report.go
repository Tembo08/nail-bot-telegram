package services

import (
	"bytes"
	"fmt"
	"nail_bot/internal/models"
	"nail_bot/internal/storage"
	"os"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type ReportService struct{}

func (s *ReportService) GetBookingsForNext3Days() ([]models.Booking, error) {
	var bookings []models.Booking
	today := time.Now().Format("2006-01-02")
	threeDaysLater := time.Now().AddDate(0, 0, 3).Format("2006-01-02")

	result := storage.GetDB().
		Where("date >= ? AND date <= ? AND status != ?", today, threeDaysLater, "cancelled").
		Order("date ASC, time ASC").
		Find(&bookings)

	return bookings, result.Error
}

func (s *ReportService) GenerateReport(bookings []models.Booking) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Путь к шрифту Liberation Sans (поддерживает кириллицу)
	fontPath := "fonts/LiberationSans-Regular.ttf"

	// Проверяем, что файл шрифта существует
	if _, err := os.Stat(fontPath); os.IsNotExist(err) {
		// Если шрифта нет — используем Helvetica (без кириллицы)
		pdf.SetFont("Helvetica", "", 16)
		pdf.Cell(40, 10, "Report for 3 days")
		pdf.Ln(12)
		pdf.SetFont("Helvetica", "", 12)
		pdf.Cell(40, 10, fmt.Sprintf("Date: %s", time.Now().Format("02.01.2006")))
	} else {
		// Регистрируем шрифт с поддержкой UTF-8
		pdf.AddUTF8Font("Liberation", "", fontPath)

		// Заголовок (русский)
		pdf.SetFont("Liberation", "", 16)
		pdf.Cell(40, 10, "Отчёт о записях на 3 дня")
		pdf.Ln(12)

		pdf.SetFont("Liberation", "", 12)
		pdf.Cell(40, 10, fmt.Sprintf("Дата отчёта: %s", time.Now().Format("02.01.2006")))
	}
	pdf.Ln(15)

	if len(bookings) == 0 {
		if _, err := os.Stat(fontPath); err == nil {
			pdf.SetFont("Liberation", "", 14)
			pdf.Cell(40, 10, "Нет активных записей на ближайшие 3 дня.")
		} else {
			pdf.SetFont("Helvetica", "", 14)
			pdf.Cell(40, 10, "No active bookings for the next 3 days.")
		}

		var buf bytes.Buffer
		err := pdf.Output(&buf)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}

	// Заголовки таблицы
	if _, err := os.Stat(fontPath); err == nil {
		pdf.SetFont("Liberation", "", 12)
	} else {
		pdf.SetFont("Helvetica", "", 12)
	}
	pdf.CellFormat(30, 8, "Дата", "1", 0, "C", false, 0, "")
	pdf.CellFormat(25, 8, "Время", "1", 0, "C", false, 0, "")
	pdf.CellFormat(35, 8, "Услуга", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 8, "Клиент", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 8, "Телефон", "1", 1, "C", false, 0, "")

	// Данные
	if _, err := os.Stat(fontPath); err == nil {
		pdf.SetFont("Liberation", "", 11)
	} else {
		pdf.SetFont("Helvetica", "", 11)
	}
	for _, b := range bookings {
		pdf.CellFormat(30, 8, b.Date, "1", 0, "C", false, 0, "")
		pdf.CellFormat(25, 8, b.Time, "1", 0, "C", false, 0, "")
		pdf.CellFormat(35, 8, b.Service, "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 8, b.ClientName, "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 8, b.Phone, "1", 1, "C", false, 0, "")
	}

	// Итого
	pdf.Ln(10)
	if _, err := os.Stat(fontPath); err == nil {
		pdf.SetFont("Liberation", "", 12)
		pdf.Cell(40, 10, fmt.Sprintf("Всего записей: %d", len(bookings)))
	} else {
		pdf.SetFont("Helvetica", "", 12)
		pdf.Cell(40, 10, fmt.Sprintf("Total bookings: %d", len(bookings)))
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
