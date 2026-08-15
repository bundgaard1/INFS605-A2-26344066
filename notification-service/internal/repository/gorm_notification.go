package repository

import (
	"context"

	"gorm.io/gorm"
	"osbourne.local/notification-service/internal/domain"
)

type GORMNotificationRepository struct {
	db *gorm.DB
}

func NewGORMNotificationRepository(db *gorm.DB) *GORMNotificationRepository {
	return &GORMNotificationRepository{db: db}
}

func (r *GORMNotificationRepository) GetUserNotifications(ctx context.Context, id string) (*[]domain.Notification, error) {

	var notifications []domain.Notification

	err := r.db.WithContext(ctx).
		Where("user_id = ?", id).
		Find(&notifications).Error

	if err != nil {
		return nil, err
	}

	return &notifications, nil
}

func (r *GORMNotificationRepository) MarkNotificationAsRead(ctx context.Context, notificationID string) (*domain.Notification, error) {
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

	return &notificationEntity, nil
}

func (r *GORMNotificationRepository) CreateNotification(ctx context.Context, notification *domain.Notification) error {

	err := r.db.WithContext(ctx).Create(&notification).Error
	if err != nil {
		return err
	}

	return nil
}
