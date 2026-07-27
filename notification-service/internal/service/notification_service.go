package service

import (
	"context"

	"osbourne.local/notification-service/internal/domain"
	"osbourne.local/notification-service/internal/repository"
)

type NotificationService struct {
	repo repository.NotificationRepository
}

func NewNotificationService(repo repository.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) GetUserNotifications(ctx context.Context, id string) (*[]domain.Notification, error) {
	return s.repo.GetUserNotifications(ctx, id)
}

func (s *NotificationService) MarkNotificationAsRead(ctx context.Context, notificationID string) (*domain.Notification, error) {
	return s.repo.MarkNotificationAsRead(ctx, notificationID)
}

func (s *NotificationService) CreateNotification(ctx context.Context, userID string, title string, message string, link string) error {

	return s.repo.CreateNotification(ctx, userID, title, message, link)
}
