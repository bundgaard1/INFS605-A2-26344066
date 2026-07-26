package server

import (
	"context"

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

func (s *NotificationServer) GetUserNotifications(ctx context.Context, req *notification.NotificationsRequest) (*notification.NotificationsResponse, error) {
	return s.notificationSvc.GetUserNotifications(ctx, req.GetUserId())
}

func (s *NotificationServer) MarkNotificationAsRead(ctx context.Context, req *notification.MarkNotificationAsReadRequest) (*notification.MarkNotificationAsReadResponse, error) {
	return s.notificationSvc.MarkNotificationAsRead(ctx, req.GetNotificationId())
}
