package services

import (
	"bytes"
	"fmt"
	"nail_bot/internal/models"
	"nail_bot/internal/storage"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type ReportService struct{}

// GetBookingsForNext3Days возвращает записи на ближайшие 3 дня
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

// GenerateReport генерирует PDF-отчёт с записями
func (s *ReportService) GenerateReport(bookings []models.Booking) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Используем стандартный шрифт
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Report for 3 days")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Date: %s", time.Now().Format("02.01.2006")))
	pdf.Ln(15)

	if len(bookings) == 0 {
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(40, 10, "No active bookings for the next 3 days.")

		// Используем bytes.Buffer для получения данных
		var buf bytes.Buffer
		err := pdf.Output(&buf)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}

	// Заголовки таблицы (на английском)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(30, 8, "Date", "1", 0, "C", false, 0, "")
	pdf.CellFormat(25, 8, "Time", "1", 0, "C", false, 0, "")
	pdf.CellFormat(35, 8, "Service", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 8, "Client", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 8, "Phone", "1", 1, "C", false, 0, "")

	// Данные
	pdf.SetFont("Arial", "", 11)
	for _, b := range bookings {
		pdf.CellFormat(30, 8, b.Date, "1", 0, "C", false, 0, "")
		pdf.CellFormat(25, 8, b.Time, "1", 0, "C", false, 0, "")
		pdf.CellFormat(35, 8, b.Service, "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 8, b.ClientName, "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 8, b.Phone, "1", 1, "C", false, 0, "")
	}

	// Итого
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Total bookings: %d", len(bookings)))

	// Используем bytes.Buffer для получения данных
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
