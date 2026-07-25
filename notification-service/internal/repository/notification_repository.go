package repository

import (
	"context"

	"osbourne.local/notification-service/gen/notification"
)

// NotificationRepository definerer alle database-operationer for notifikationer
type NotificationRepository interface {
	GetUserNotifications(ctx context.Context, id string) (*notification.NotificationsResponse, error)
	MarkNotificationAsRead(ctx context.Context, notificationID string) (*notification.MarkNotificationAsReadResponse, error)
	CreateNotification(ctx context.Context, userID string, title string, message string, link string) error
}
