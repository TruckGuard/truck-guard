package main

import (
	"log/slog"

	"github.com/truckguard/core/src/models"
	"github.com/truckguard/core/src/repository"
)

func SeedData() {
	slog.Info("Seeding system settings...")

	settings := []models.SystemSetting{
		{
			Key:         "permit_timeout_seconds",
			Name:        "Таймаут перепустки (сек)",
			Description: "Час, протягом якого перепустка залишається активною для додавання подій (камери/ваги)",
			Value:       "60",
			Default:     "60",
		},
		{
			Key:         "unverified_permits_cleanup_hours",
			Name:        "Очищення неверіфікованих перепусток (год)",
			Description: "Час, після якого неверіфіковані перепустки стають недійсними (is_void = true)",
			Value:       "24",
			Default:     "24",
		},
		{
			Key:         "permit_audit_cleanup_days",
			Name:        "Очищення аудиту (днів)",
			Description: "Через скільки днів видаляти логи аудиту для закритих перепусток",
			Value:       "7",
			Default:     "7",
		},
		{
			Key:         "void_permits_deletion_days",
			Name:        "Видалення анульованих (днів)",
			Description: "Через скільки днів остаточно видаляти анульовані перепустки",
			Value:       "30",
			Default:     "30",
		},
		{
			Key:         "fuzzy_match_max_distance",
			Name:        "Поріг нечіткого збігу (Левенштейн)",
			Description: "Максимальна відстань Левенштейна для нечіткого пошуку номерів. 0 = вимкнено.",
			Value:       "2",
			Default:     "2",
		},
		{
			Key:         "notification_retention_days",
			Name:        "Зберігання сповіщень (днів)",
			Description: "Через скільки днів автоматично видаляти старі сповіщення",
			Value:       "2",
			Default:     "2",
		},
	}

	for _, s := range settings {
		var existing models.SystemSetting
		if err := repository.DB.Where("key = ?", s.Key).First(&existing).Error; err != nil {
			if err := repository.DB.Create(&s).Error; err != nil {
				slog.Error("Failed to seed setting", "key", s.Key, "error", err)
			} else {
				slog.Info("Created system setting", "key", s.Key)
			}
		} else {
			// Optionally update fields other than Value if they changed in code
			existing.Name = s.Name
			existing.Description = s.Description
			existing.Default = s.Default
			repository.DB.Save(&existing)
		}
	}

	slog.Info("System settings seeding completed")
}
