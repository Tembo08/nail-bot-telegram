package services

import (
	"nail_bot/internal/models"
	"nail_bot/internal/storage"
	"time"
)

type SessionService struct{}

func (s *SessionService) GetOrCreateSession(userID int64) (*models.UserSession, error) {
	var session models.UserSession
	result := storage.GetDB().Where("user_id = ?", userID).First(&session)

	if result.Error != nil {
		// Сессии нет — создаём новую
		session = models.UserSession{
			UserID:    userID,
			Step:      "service",
			UpdatedAt: time.Now(),
		}
		if err := storage.GetDB().Create(&session).Error; err != nil {
			return nil, err
		}
		return &session, nil
	}

	return &session, nil
}

func (s *SessionService) UpdateSession(userID int64, step string, data map[string]string) error {
	var session models.UserSession
	result := storage.GetDB().Where("user_id = ?", userID).First(&session)
	if result.Error != nil {
		return result.Error
	}

	session.Step = step
	session.UpdatedAt = time.Now()

	if service, ok := data["service"]; ok {
		session.Service = service
	}
	if date, ok := data["date"]; ok {
		session.Date = date
	}
	if timeStr, ok := data["time"]; ok {
		session.Time = timeStr
	}
	if contact, ok := data["contact"]; ok {
		session.Contact = contact
	}

	return storage.GetDB().Save(&session).Error
}

func (s *SessionService) DeleteSession(userID int64) error {
	return storage.GetDB().Where("user_id = ?", userID).Delete(&models.UserSession{}).Error
}

func (s *SessionService) GetSession(userID int64) (*models.UserSession, error) {
	var session models.UserSession
	result := storage.GetDB().Where("user_id = ?", userID).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}
