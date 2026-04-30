package repository

import (
	"context"
	"time"

	"github.com/truckguard/core/src/models"
)

// CleanupNotifications deletes notifications older than the specified number of days.
func CleanupNotifications(ctx context.Context, days int) (int64, error) {
	threshold := time.Now().Add(-time.Duration(days*24) * time.Hour)
	result := DB.WithContext(ctx).Unscoped().
		Where("created_at < ?", threshold).
		Delete(&models.Notification{})
	return result.RowsAffected, result.Error
}

// GetNotifications returns a list of notifications for a post, optionally filtering by unread status.
func GetNotifications(ctx context.Context, postID uint, limit, offset int, unreadOnly bool) ([]models.Notification, error) {
	var notifs []models.Notification
	query := DB.WithContext(ctx).Where("post_id = ?", postID)
	if unreadOnly {
		query = query.Where("is_read = ?", false)
	}
	err := query.Order("created_at desc").Limit(limit).Offset(offset).Find(&notifs).Error
	return notifs, err
}

// MarkNotificationRead marks a specific notification as read.
func MarkNotificationRead(ctx context.Context, id uint) error {
	return DB.WithContext(ctx).Model(&models.Notification{}).Where("id = ?", id).Update("is_read", true).Error
}

// MarkAllNotificationsRead marks all notifications as read for a given post.
func MarkAllNotificationsRead(ctx context.Context, postID uint) error {
	return DB.WithContext(ctx).Model(&models.Notification{}).Where("post_id = ? AND is_read = ?", postID, false).Update("is_read", true).Error
}
