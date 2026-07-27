package repository

import (
	"context"

	"osbourne.local/notification-service/internal/domain"
)

// NotificationRepository definerer alle database-operationer for notifikationer
type NotificationRepository interface {
	GetUserNotifications(ctx context.Context, id string) (*[]domain.Notification, error)
	MarkNotificationAsRead(ctx context.Context, notificationID string) (*domain.Notification, error)
	CreateNotification(ctx context.Context, userID string, title string, message string, link string) error
}
