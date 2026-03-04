package datarepo

import (
	"context"

	"github.com/truckguard/core/src/models"
	"github.com/truckguard/core/src/repository"
)

type PostFilter struct {
	Name string
}

func ListPosts(ctx context.Context, limit, offset int, filter PostFilter) ([]models.CustomsPost, int64, error) {
	var posts []models.CustomsPost
	var total int64

	query := repository.DB.WithContext(ctx).Model(&models.CustomsPost{})
	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Limit(limit).Offset(offset).Order("id desc").Find(&posts).Error; err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}

func GetPost(ctx context.Context, id string) (models.CustomsPost, error) {
	var post models.CustomsPost
	err := repository.DB.WithContext(ctx).First(&post, id).Error
	return post, err
}

func CreatePost(ctx context.Context, input *models.CustomsPost) error {
	return repository.DB.WithContext(ctx).Create(input).Error
}

func UpdatePost(ctx context.Context, id string, input *models.CustomsPost) (models.CustomsPost, error) {
	post, err := GetPost(ctx, id)
	if err != nil {
		return post, err
	}
	post.Name = input.Name
	post.Description = input.Description
	if err := repository.DB.WithContext(ctx).Save(&post).Error; err != nil {
		return post, err
	}
	return post, nil
}

func DeletePost(ctx context.Context, id string) error {
	return repository.DB.WithContext(ctx).Unscoped().Delete(&models.CustomsPost{}, id).Error
}
