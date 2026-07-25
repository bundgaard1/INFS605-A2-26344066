package service

import (
	"context"

	"osbourne.local/notification-service/gen/notification"
	"osbourne.local/notification-service/internal/repository"
)

type NotificationService struct {
	repo repository.NotificationRepository
}

func NewNotificationService(repo repository.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) GetUserNotifications(ctx context.Context, id string) (*notification.NotificationsResponse, error) {
	return s.repo.GetUserNotifications(ctx, id)
}

func (s *NotificationService) MarkNotificationAsRead(ctx context.Context, notificationID string) (*notification.MarkNotificationAsReadResponse, error) {
	return s.repo.MarkNotificationAsRead(ctx, notificationID)
}

func (s *NotificationService) CreateNotification(ctx context.Context, userID string, title string, message string, link string) error {

	return s.repo.CreateNotification(ctx, userID, title, message, link)
}
