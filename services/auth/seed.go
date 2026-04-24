package main

import (
	"os"

	"github.com/truckguard/auth/src/models"
	"github.com/truckguard/auth/src/repository"
	"golang.org/x/crypto/bcrypt"
)

func genCRUDPerms(resource, name, module string, includeAll bool) []models.Permission {
	perms := []models.Permission{
		{ID: "read:" + resource, Name: name + ": Читання", Module: module},
		{ID: "create:" + resource, Name: name + ": Створення", Module: module},
		{ID: "update:" + resource, Name: name + ": Оновлення", Module: module},
		{ID: "delete:" + resource, Name: name + ": Видалення", Module: module},
		{ID: "manage:" + resource, Name: name + ": Повний доступ", Module: module},
	}
	if includeAll {
		perms = append(perms, []models.Permission{
			{ID: "read:" + resource + ":all", Name: name + ": Читання (всі)", Module: module},
			{ID: "manage:" + resource + ":all", Name: name + ": Повний доступ (всі)", Module: module},
		}...)
	}
	return perms
}

func seedData() {
	var perms []models.Permission

	// Auth
	perms = append(perms, genCRUDPerms("users", "Користувачі", "auth", true)...)
	perms = append(perms, genCRUDPerms("roles", "Ролі", "auth", false)...)
	perms = append(perms, genCRUDPerms("keys", "API Ключі", "auth", false)...)
	perms = append(perms, models.Permission{ID: "read:audit", Name: "Аудит: Перегляд", Module: "auth"})

	// Ingest
	perms = append(perms, models.Permission{ID: "ingest:events", Name: "Імпорт: Створення", Module: "ingestor"})

	// Core
	perms = append(perms, genCRUDPerms("cameras", "Камери", "cameras", true)...)
	perms = append(perms, genCRUDPerms("scales", "Ваги", "scales", true)...)
	perms = append(perms, genCRUDPerms("events", "Події", "events", true)...)

	permitsPerms := genCRUDPerms("permits", "Перепустки", "permits", true)
	permitsPerms = append(permitsPerms, []models.Permission{
		{ID: "validate:permits", Name: "Перепустки: Валідація", Module: "permits"},
		{ID: "validate:permits:all", Name: "Перепустки: Валідація всіх", Module: "permits"},
	}...)
	perms = append(perms, permitsPerms...)

	perms = append(perms, genCRUDPerms("settings", "Налаштування", "settings", false)...)
	perms = append(perms, genCRUDPerms("data", "Довідники", "data", false)...)

	// Customs Parser
	perms = append(perms, models.Permission{ID: "read:customs", Name: "Митниця: Отримання даних", Module: "customs-parser"})
	perms = append(perms, models.Permission{ID: "read:companies", Name: "Компанії: Отримання даних", Module: "customs-parser"})

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
	var managerPerms []models.Permission
	repository.DB.Where("id IN ?", managerPermIDs).Find(&managerPerms)
	repository.DB.Model(&managerRole).Association("Permissions").Replace(managerPerms)

	// 3. Оператор - тільки перегляд та створення подій/перепусток
	var operatorRole models.Role
	repository.DB.FirstOrCreate(&operatorRole, models.Role{Name: "operator", Description: "Оператор (перегляд та реєстрація)"})
	operatorPermIDs := []string{
		"read:permits", "create:permits", "manage:permits",
		"read:events", "create:events", "manage:events",
		"read:cameras", "read:customs",
	}
	var operatorPerms []models.Permission
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
			var workerPerms []models.Permission
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
		{Method: "POST", PathPattern: `^/auth/admin/users/.*/reset-password$`, RequiredPermission: "manage:users", Description: "Скидання пароля користувача"},
		{Method: "GET", PathPattern: `^/auth/admin/roles.*`, RequiredPermission: "read:roles", Description: "Перегляд ролей"},
		{Method: "POST", PathPattern: `^/auth/admin/roles.*`, RequiredPermission: "manage:roles", Description: "Створення ролей"},
		{Method: "CRUD", PathPattern: `^/auth/admin/roles/.*`, RequiredPermission: "roles", Description: "Керування ролями"},
		{Method: "GET", PathPattern: `^/auth/admin/keys.*`, RequiredPermission: "read:keys", Description: "Перегляд ключів"},
		{Method: "POST", PathPattern: `^/auth/admin/keys.*`, RequiredPermission: "manage:keys", Description: "Створення ключів"},
		{Method: "CRUD", PathPattern: `^/auth/admin/keys/.*`, RequiredPermission: "keys", Description: "Керування ключами"},
		{Method: "GET", PathPattern: `^/auth/admin/permissions$`, RequiredPermission: "", Description: "Список всіх дозволів"},

		// Ingestor Service
		{Method: "POST", PathPattern: `^/ingest/.*`, RequiredPermission: "ingest:events", Description: "Імпорт даних"},

		// Core Service: CONFIGS
		{Method: "CRUD", PathPattern: `^/api/configs/cameras.*`, RequiredPermission: "cameras,cameras:all", Description: "Керування камерами"},
		{Method: "CRUD", PathPattern: `^/api/configs/scales.*`, RequiredPermission: "scales,scales:all", Description: "Керування вагами"},
		{Method: "CRUD", PathPattern: `^/api/configs/settings.*`, RequiredPermission: "settings", Description: "Керування налаштуваннями"},
		{Method: "CRUD", PathPattern: `^/api/configs/excluded-plates.*`, RequiredPermission: "settings", Description: "Керування чорним списком"},

		// Core Service: DATA
		{Method: "CRUD", PathPattern: `^/api/data/.*`, RequiredPermission: "data", Description: "Керування довідниками"},

		// Core Service: EVENTS
		{Method: "CRUD", PathPattern: `^/api/events/.*`, RequiredPermission: "events,events:all", Description: "Керування подіями"},

		// Core Service: PERMITS
		{Method: "CRUD", PathPattern: `^/api/permits.*`, RequiredPermission: "permits,permits:all", Description: "Керування перепустками"},
		{Method: "POST", PathPattern: `^/api/permits/.*/validate$`, RequiredPermission: "validate:permits,validate:permits:all", Description: "Валідація перепусток"},

		// Core Service: USERS
		{Method: "GET", PathPattern: `^/api/users/?$`, RequiredPermission: "read:users", Description: "Список профілів"},
		{Method: "GET", PathPattern: `^/api/users/me`, RequiredPermission: "", Description: "Мій профіль"},
		{Method: "PUT", PathPattern: `^/api/users/me`, RequiredPermission: "", Description: "Оновлення мого профілю"},
		{Method: "POST", PathPattern: `^/api/users/?$`, RequiredPermission: "manage:users,create:users", Description: "Створення профілю"},
		{Method: "CRUD", PathPattern: `^/api/users/.*`, RequiredPermission: "users", Description: "Керування профілями"},

		// Images
		{Method: "GET", PathPattern: `^/api/images/.*`, RequiredPermission: "read:events,read:events:all", Description: "Перегляд зображень"},

		// Customs
		{Method: "GET", PathPattern: `^/api/data-parser/customs/.*`, RequiredPermission: "read:customs", Description: "Доступ до даних з митниці"},

		// EDR / Company registry
		{Method: "GET", PathPattern: `^/api/data-parser/companies/.*`, RequiredPermission: "read:companies", Description: "Пошук компаній у реєстрі ЄДР"},

		// Audit
		{Method: "GET", PathPattern: `^/api/audit.*`, RequiredPermission: "read:audit", Description: "Журнал аудиту"},
	}

	for _, r := range rules {
		var existing models.PolicyRule
		if err := repository.DB.Where("method = ? AND path_pattern = ?", r.Method, r.PathPattern).First(&existing).Error; err != nil {
			repository.DB.Create(&r)
		} else {
			existing.RequiredPermission = r.RequiredPermission
			repository.DB.Save(&existing)
		}
	}

	// Permission Hierarchy
	hierarchy := []models.PermissionHierarchy{
		{ParentID: "manage:users", ChildID: "read:users"},
		{ParentID: "manage:users", ChildID: "create:users"},
		{ParentID: "manage:users", ChildID: "update:users"},
		{ParentID: "manage:users", ChildID: "delete:users"},
		{ParentID: "manage:users", ChildID: "read:roles"},
		{ParentID: "manage:users", ChildID: "read:data"},

		{ParentID: "manage:roles", ChildID: "read:roles"},
		{ParentID: "manage:roles", ChildID: "create:roles"},
		{ParentID: "manage:roles", ChildID: "update:roles"},
		{ParentID: "manage:roles", ChildID: "delete:roles"},
		{ParentID: "manage:roles", ChildID: "read:data"},

		{ParentID: "manage:permits", ChildID: "read:permits"},
		{ParentID: "manage:permits", ChildID: "create:permits"},
		{ParentID: "manage:permits", ChildID: "update:permits"},
		{ParentID: "manage:permits", ChildID: "delete:permits"},
		{ParentID: "manage:permits", ChildID: "validate:permits"},
		{ParentID: "manage:permits", ChildID: "read:data"},
		{ParentID: "manage:permits", ChildID: "read:events"},
		{ParentID: "manage:permits", ChildID: "read:customs"},

		{ParentID: "read:permits", ChildID: "read:data"},
		{ParentID: "read:permits", ChildID: "read:events"},

		{ParentID: "manage:events", ChildID: "read:events"},
		{ParentID: "manage:events", ChildID: "create:events"},
		{ParentID: "manage:events", ChildID: "update:events"},
		{ParentID: "manage:events", ChildID: "delete:events"},
		{ParentID: "manage:events", ChildID: "read:cameras"},
		{ParentID: "manage:events", ChildID: "read:data"},

		{ParentID: "manage:cameras", ChildID: "read:cameras"},
		{ParentID: "manage:cameras", ChildID: "create:cameras"},
		{ParentID: "manage:cameras", ChildID: "update:cameras"},
		{ParentID: "manage:cameras", ChildID: "delete:cameras"},
		{ParentID: "manage:cameras", ChildID: "read:data"},

		{ParentID: "read:cameras", ChildID: "read:data"},

		{ParentID: "manage:scales", ChildID: "read:scales"},
		{ParentID: "manage:scales", ChildID: "create:scales"},
		{ParentID: "manage:scales", ChildID: "update:scales"},
		{ParentID: "manage:scales", ChildID: "delete:scales"},
		{ParentID: "manage:scales", ChildID: "read:data"},

		{ParentID: "read:scales", ChildID: "read:data"},

		{ParentID: "manage:settings", ChildID: "read:settings"},
		{ParentID: "manage:settings", ChildID: "create:settings"},
		{ParentID: "manage:settings", ChildID: "update:settings"},
		{ParentID: "manage:settings", ChildID: "delete:settings"},
		{ParentID: "manage:settings", ChildID: "read:data"},
	}

	for _, h := range hierarchy {
		repository.DB.Where("parent_id = ? AND child_id = ?", h.ParentID, h.ChildID).FirstOrCreate(&h)
	}
}
