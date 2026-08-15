package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"osbourne.local/notification-service/internal/domain"
)

type NotificationService struct {
	repo domain.NotificationRepository
}

func NewNotificationService(repo domain.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) GetUserNotifications(ctx context.Context, id string) (*[]domain.Notification, error) {
	return s.repo.GetUserNotifications(ctx, id)
}

func (s *NotificationService) MarkNotificationAsRead(ctx context.Context, notificationID string) (*domain.Notification, error) {
	return s.repo.MarkNotificationAsRead(ctx, notificationID)
}

func (s *NotificationService) CreateNotification(ctx context.Context, userID string, title string, message string, link string) error {
	n := &domain.Notification{
		ID:        uuid.NewString(),
		UserID:    userID,
		Title:     title,
		Message:   message,
		LinkURL:   link,
		IsRead:    false,
		CreatedAt: time.Now(),
	}
	return s.repo.CreateNotification(ctx, n)
}
