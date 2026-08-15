package domain

import (
	"context"
)

// NotificationRepository definerer alle database-operationer for notifikationer
type NotificationRepository interface {
	GetUserNotifications(context.Context, string) (*[]Notification, error)
	MarkNotificationAsRead(context.Context, string) (*Notification, error)
	CreateNotification(context.Context, *Notification) error
}
