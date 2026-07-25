package server

import (
	"osbourne.local/notification-service/gen/notification"
	"osbourne.local/notification-service/internal/service"
)

type NotificationServer struct {
	notification.UnimplementedNotificationServiceServer
	notificationSvc *service.NotificationService
}

func NewNotificationServer(notificationSvc *service.NotificationService) *NotificationServer {
	return &NotificationServer{
		notificationSvc: notificationSvc,
	}
}
