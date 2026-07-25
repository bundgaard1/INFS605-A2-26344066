package repository

import (
	"context"

	"gorm.io/gorm"
	"osbourne.local/notification-service/gen/notification"
	"osbourne.local/notification-service/internal/domain"
)

type GORMNotificationRepository struct {
	db *gorm.DB
}

func NewGORMNotificationRepository(db *gorm.DB) *GORMNotificationRepository {
	return &GORMNotificationRepository{db: db}
}

func (r *GORMNotificationRepository) GetUserNotifications(ctx context.Context, id string) (*notification.NotificationsResponse, error) {

	var notifications []domain.Notification

	err := r.db.WithContext(ctx).
		Where("user_id = ?", id).
		Find(&notifications).Error

	if err != nil {
		return nil, err
	}

	var notificationResponses []*notification.Notification

	for _, n := range notifications {
		notificationResponses = append(notificationResponses, &notification.Notification{
			Id:        n.ID,
			UserId:    n.UserID,
			Title:     n.Title,
			Msg:       n.Message,
			Link:      n.LinkURL,
			IsRead:    n.IsRead,
			Timestamp: n.CreatedAt.Unix(),
		})
	}

	return &notification.NotificationsResponse{
		Notifications: notificationResponses,
	}, nil
}

func (r *GORMNotificationRepository) MarkNotificationAsRead(ctx context.Context, notificationID string) (*notification.MarkNotificationAsReadResponse, error) {
	var notificationEntity domain.Notification

	err := r.db.WithContext(ctx).
		First(&notificationEntity, "id = ?", notificationID).Error

	if err != nil {
		return nil, err
	}

	notificationEntity.IsRead = true

	err = r.db.WithContext(ctx).Save(&notificationEntity).Error
	if err != nil {
		return nil, err
	}

	return &notification.MarkNotificationAsReadResponse{
		Success: true,
	}, nil
}

func (r *GORMNotificationRepository) CreateNotification(ctx context.Context, userID string, title string, message string, link string) error {
	notification := domain.Notification{
		UserID:  userID,
		Title:   title,
		Message: message,
		LinkURL: link,
	}

	err := r.db.WithContext(ctx).Create(&notification).Error
	if err != nil {
		return err
	}

	return nil
}
