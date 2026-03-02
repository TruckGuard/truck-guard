package main

import (
	"os"

	"github.com/truckguard/auth/src/models"
	"github.com/truckguard/auth/src/repository"
	"golang.org/x/crypto/bcrypt"
)

func seedData() {
	perms := []models.Permission{
		// Auth
		{ID: "read:users", Name: "Користувачі: Перегляд", Module: "auth"},
		{ID: "read:users:all", Name: "Користувачі: Перегляд всіх", Module: "auth"},
		{ID: "manage:users", Name: "Користувачі: Повний доступ", Module: "auth"},
		{ID: "manage:users:all", Name: "Користувачі: Повний доступ (всі)", Module: "auth"},

		{ID: "read:roles", Name: "Ролі: Перегляд", Module: "auth"},
		{ID: "manage:roles", Name: "Ролі: Повний доступ", Module: "auth"},

		{ID: "read:keys", Name: "API Ключі: Перегляд", Module: "auth"},
		{ID: "manage:keys", Name: "API Ключі: Повний доступ", Module: "auth"},

		// Ingest
		{ID: "manage:ingest", Name: "Імпорт: Повний доступ", Module: "ingestor"},

		// Core
		{ID: "read:cameras", Name: "Камери: Перегляд", Module: "core"},
		{ID: "read:cameras:all", Name: "Камери: Перегляд всіх", Module: "core"},
		{ID: "manage:cameras", Name: "Камери: Редагування", Module: "core"},
		{ID: "manage:cameras:all", Name: "Камери: Редагування всіх", Module: "core"},

		{ID: "read:scales", Name: "Ваги: Перегляд", Module: "core"},
		{ID: "read:scales:all", Name: "Ваги: Перегляд всіх", Module: "core"},
		{ID: "manage:scales", Name: "Ваги: Редагування", Module: "core"},
		{ID: "manage:scales:all", Name: "Ваги: Редагування всіх", Module: "core"},

		{ID: "read:events", Name: "Події: Перегляд", Module: "core"},
		{ID: "read:events:all", Name: "Події: Перегляд всіх", Module: "core"},
		{ID: "manage:events", Name: "Події: Повний доступ", Module: "core"},
		{ID: "manage:events:all", Name: "Події: Повний доступ (всі)", Module: "core"},

		{ID: "read:permits", Name: "Перепустки: Перегляд", Module: "core"},
		{ID: "read:permits:all", Name: "Перепустки: Перегляд всіх", Module: "core"},
		{ID: "validate:permits", Name: "Перепустки: Валідація", Module: "core"},
		{ID: "validate:permits:all", Name: "Перепустки: Валідація всіх", Module: "core"},
		{ID: "manage:permits", Name: "Перепустки: Повний доступ", Module: "core"},
		{ID: "manage:permits:all", Name: "Перепустки: Повний доступ (всі)", Module: "core"},

		{ID: "read:settings", Name: "Налаштування: Перегляд", Module: "core"},
		{ID: "manage:settings", Name: "Налаштування: Повний доступ", Module: "core"},

		{ID: "read:data", Name: "Довідники: Перегляд", Module: "core"},
		{ID: "manage:data", Name: "Довідники: Повний доступ", Module: "core"},

		{ID: "read:audit", Name: "Аудит: Перегляд", Module: "auth"},

		// Customs Parser
		{ID: "read:customs", Name: "Митниця: Отримання даних", Module: "customs-parser"},
	}

	for _, p := range perms {
		repository.DB.Save(&p)
	}

	// 1. Адмін - має все
	var adminRole models.Role
	repository.DB.FirstOrCreate(&adminRole, models.Role{Name: "admin", Description: "Повний доступ до системи"})
	repository.DB.Model(&adminRole).Association("Permissions").Replace(perms)

	// 2. Менеджер - може редагувати, але не видаляти
	var managerRole models.Role
	repository.DB.FirstOrCreate(&managerRole, models.Role{Name: "manager", Description: "Керівник (редагування даних)"})
	managerPermIDs := []string{
		"manage:users", "manage:permits", "validate:permits", "manage:events", "manage:settings",
		"manage:cameras", "manage:scales", "read:roles", "read:keys", "read:customs", "read:audit",
	}
	managerPerms := []models.Permission{}
	repository.DB.Where("id IN ?", managerPermIDs).Find(&managerPerms)
	repository.DB.Model(&managerRole).Association("Permissions").Replace(managerPerms)

	// 3. Оператор - тільки перегляд та створення подій/перепусток
	var operatorRole models.Role
	repository.DB.FirstOrCreate(&operatorRole, models.Role{Name: "operator", Description: "Оператор (перегляд та реєстрація)"})
	operatorPermIDs := []string{
		"read:permits", "manage:permits", "read:events", "manage:events", "read:cameras", "read:customs",
	}
	// Примітка: для оператора manage може бути обмежений core-сервісом,
	// але тут ми даємо базовий дозвіл на ці ресурси.
	operatorPerms := []models.Permission{}
	repository.DB.Where("id IN ?", operatorPermIDs).Find(&operatorPerms)
	repository.DB.Model(&operatorRole).Association("Permissions").Replace(operatorPerms)

	adminUsername := "admin"
	adminPassword := os.Getenv("ADMIN_DEFAULT_PASSWORD")
	if adminPassword == "" {
		adminPassword = "admin123"
	}
	var adminUser models.User
	err := repository.DB.Where("username = ?", adminUsername).First(&adminUser).Error
	if err != nil {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)

		newAdmin := models.User{
			Username:     adminUsername,
			PasswordHash: string(hashedPassword),
			RoleID:       adminRole.ID,
			Role:         adminRole,
		}

		if createErr := repository.DB.Create(&newAdmin).Error; createErr == nil {
			println("Адміністратора за замовчуванням успішно створено: admin")
		}
	}

	workerKey := os.Getenv("WORKER_SYSTEM_KEY")
	if workerKey != "" {
		h := repository.HashKey(workerKey)
		var existingKey models.APIKey
		err := repository.DB.Where("key_hash = ?", h).First(&existingKey).Error
		if err != nil {
			workerPerms := []models.Permission{}
			repository.DB.Where("id IN ?", []string{"manage:events", "read:cameras:all", "read:scales:all"}).Find(&workerPerms)

			newKey := models.APIKey{
				KeyHash:     h,
				OwnerName:   "Системний воркер",
				IsActive:    true,
				Permissions: workerPerms,
			}
			repository.DB.Create(&newKey)
			println("API ключ для системного воркера успішно додано")
		}
	}

	rules := []models.PolicyRule{
		// Auth Service
		{Method: "POST", PathPattern: `^/auth/register$`, RequiredPermission: "manage:users", Description: "Реєстрація користувачів"},
		{Method: "GET", PathPattern: `^/auth/admin/users.*`, RequiredPermission: "read:users", Description: "Перегляд користувачів"},
		{Method: "PUT", PathPattern: `^/auth/admin/users/.*/role$`, RequiredPermission: "manage:users", Description: "Зміна ролі користувача"},
		{Method: "DELETE", PathPattern: `^/auth/admin/users/.*`, RequiredPermission: "manage:users", Description: "Видалення користувача"},
		{Method: "GET", PathPattern: `^/auth/admin/roles.*`, RequiredPermission: "read:roles", Description: "Перегляд ролей"},
		{Method: "POST", PathPattern: `^/auth/admin/roles.*`, RequiredPermission: "manage:roles", Description: "Створення ролей"},
		{Method: "*", PathPattern: `^/auth/admin/roles/.*`, RequiredPermission: "manage:roles", Description: "Керування ролями"},
		{Method: "GET", PathPattern: `^/auth/admin/keys.*`, RequiredPermission: "read:keys", Description: "Перегляд ключів"},
		{Method: "POST", PathPattern: `^/auth/admin/keys.*`, RequiredPermission: "manage:keys", Description: "Створення ключів"},
		{Method: "*", PathPattern: `^/auth/admin/keys/.*`, RequiredPermission: "manage:keys", Description: "Керування ключами"},
		{Method: "GET", PathPattern: `^/auth/admin/permissions$`, RequiredPermission: "read:roles", Description: "Список всіх дозволів"},

		// Ingestor Service
		{Method: "POST", PathPattern: `^/ingest/.*`, RequiredPermission: "manage:ingest", Description: "Імпорт даних"},

		// Core Service: CONFIGS
		{Method: "GET", PathPattern: `^/api/configs/cameras.*`, RequiredPermission: "read:cameras", Description: "Перегляд камер"},
		{Method: "GET", PathPattern: `^/api/configs/cameras/.*`, RequiredPermission: "read:cameras:all", Description: "Перегляд всіх камер"},
		{Method: "*", PathPattern: `^/api/configs/cameras.*`, RequiredPermission: "manage:cameras", Description: "Керування камерами"},
		{Method: "*", PathPattern: `^/api/configs/cameras.*`, RequiredPermission: "manage:cameras:all", Description: "Керування всіма камерами"},

		{Method: "GET", PathPattern: `^/api/configs/scales.*`, RequiredPermission: "read:scales", Description: "Перегляд ваг"},
		{Method: "GET", PathPattern: `^/api/configs/scales/.*`, RequiredPermission: "read:scales:all", Description: "Перегляд всіх ваг"},
		{Method: "*", PathPattern: `^/api/configs/scales.*`, RequiredPermission: "manage:scales", Description: "Керування вагами"},
		{Method: "*", PathPattern: `^/api/configs/scales.*`, RequiredPermission: "manage:scales:all", Description: "Керування всіма вагами"},

		{Method: "GET", PathPattern: `^/api/configs/settings.*`, RequiredPermission: "read:settings", Description: "Перегляд налаштувань"},
		{Method: "*", PathPattern: `^/api/configs/settings.*`, RequiredPermission: "manage:settings", Description: "Керування налаштуваннями"},

		{Method: "GET", PathPattern: `^/api/configs/excluded-plates.*`, RequiredPermission: "read:settings", Description: "Перегляд чорного списку"},
		{Method: "*", PathPattern: `^/api/configs/excluded-plates.*`, RequiredPermission: "manage:settings", Description: "Керування чорним списком"},

		// Core Service: DATA
		{Method: "GET", PathPattern: `^/api/data/.*`, RequiredPermission: "read:data", Description: "Перегляд довідників"},
		{Method: "*", PathPattern: `^/api/data/.*`, RequiredPermission: "manage:data", Description: "Керування довідниками"},

		// Core Service: EVENTS
		{Method: "GET", PathPattern: `^/api/events/.*`, RequiredPermission: "read:events", Description: "Перегляд подій"},
		{Method: "GET", PathPattern: `^/api/events/.*`, RequiredPermission: "read:events:all", Description: "Перегляд всіх подій"},
		{Method: "POST", PathPattern: `^/api/events/.*`, RequiredPermission: "manage:events", Description: "Реєстрація подій"},
		{Method: "PATCH", PathPattern: `^/api/events/plate/.*`, RequiredPermission: "manage:events", Description: "Корекція номерів"},
		{Method: "PATCH", PathPattern: `^/api/events/plate/.*`, RequiredPermission: "manage:events:all", Description: "Корекція всіх номерів"},

		// Core Service: PERMITS
		{Method: "GET", PathPattern: `^/api/permits.*`, RequiredPermission: "read:permits", Description: "Перегляд перепусток"},
		{Method: "POST", PathPattern: `^/api/permits/.*/validate$`, RequiredPermission: "validate:permits", Description: "Валідація перепусток"},
		{Method: "POST", PathPattern: `^/api/permits/.*/validate$`, RequiredPermission: "validate:permits:all", Description: "Валідація всіх перепусток"},
		{Method: "*", PathPattern: `^/api/permits.*`, RequiredPermission: "manage:permits", Description: "Керування перепустками"},
		{Method: "*", PathPattern: `^/api/permits.*`, RequiredPermission: "manage:permits:all", Description: "Керування всіма перепустками"},

		// Core Service: USERS
		{Method: "GET", PathPattern: `^/api/users/?$`, RequiredPermission: "read:users", Description: "Список профілів"},
		{Method: "GET", PathPattern: `^/api/users/me`, RequiredPermission: "", Description: "Мій профіль"},
		{Method: "PUT", PathPattern: `^/api/users/me`, RequiredPermission: "", Description: "Оновлення мого профілю"},
		{Method: "POST", PathPattern: `^/api/users/?$`, RequiredPermission: "manage:users", Description: "Створення профілю"},
		{Method: "*", PathPattern: `^/api/users/.*`, RequiredPermission: "manage:users", Description: "Керування профілями"},

		// Images
		{Method: "GET", PathPattern: `^/api/images/.*`, RequiredPermission: "read:events", Description: "Перегляд зображень"},
		{Method: "GET", PathPattern: `^/api/images/.*`, RequiredPermission: "read:events:all", Description: "Перегляд зображень"},

		// Customs
		{Method: "GET", PathPattern: `^/api/data-parser/customs/.*`, RequiredPermission: "read:customs", Description: "Доступ до даних з митниці"},
	}

	for _, r := range rules {
		var existing models.PolicyRule
		if err := repository.DB.Where("method = ? AND path_pattern = ?", r.Method, r.PathPattern).First(&existing).Error; err != nil {
			repository.DB.Create(&r)
		}
	}
}
