package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/truckguard/core/src/models"
)

func GetUserPostIDCached(ctx context.Context, authID string) (string, error) {
	cacheKey := fmt.Sprintf("user:%s:post_id", authID)

	// Try Redis first
	val, err := RDB.Get(ctx, cacheKey).Result()
	if err == nil {
		if val == "nil" {
			return "", nil // explicitly null
		}
		return val, nil
	}

	// Fallback to DB
	var user models.User
	if err := DB.WithContext(ctx).Where("auth_id = ?", authID).First(&user).Error; err != nil {
		return "", err
	}

	postIDStr := "nil"
	if user.CustomsPostID != nil {
		postIDStr = fmt.Sprintf("%d", *user.CustomsPostID)
	}

	// Cache for 30 minutes
	RDB.Set(ctx, cacheKey, postIDStr, 30*time.Minute)

	if postIDStr == "nil" {
		return "", nil
	}

	return postIDStr, nil
}

func GetUserByAuthID(ctx context.Context, authID uint) (models.User, error) {
	var user models.User
	err := DB.WithContext(ctx).Where("auth_id = ?", authID).First(&user).Error
	return user, err
}

func UpdateUser(ctx context.Context, user *models.User) error {
	var err error
	if user.ID == 0 {
		err = DB.WithContext(ctx).Create(user).Error
	} else {
		err = DB.WithContext(ctx).Save(user).Error
	}

	if err == nil {
		InvalidateUserCache(ctx, fmt.Sprintf("%d", user.AuthID))
	}
	return err
}

func InvalidateUserCache(ctx context.Context, authID string) {
	cacheKey := fmt.Sprintf("user:%s:post_id", authID)
	RDB.Del(ctx, cacheKey)
}

func DeleteUser(ctx context.Context, authID uint) error {
	err := DB.WithContext(ctx).Where("auth_id = ?", authID).Delete(&models.User{}).Error
	if err == nil {
		InvalidateUserCache(ctx, fmt.Sprintf("%d", authID))
	}
	return err
}
