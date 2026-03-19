package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/truckguard/auth/src/models"
)
var PermissionHierarchy = make(map[string][]string)

// LoadPermissionHierarchy завантажує ієрархію дозволів з бази даних.
func LoadPermissionHierarchy() {
	var relations []models.PermissionHierarchy
	if err := DB.Find(&relations).Error; err != nil {
		slog.Error("Failed to load permission hierarchy from DB", "error", err)
		return
	}

	newHierarchy := make(map[string][]string)
	for _, r := range relations {
		newHierarchy[r.ParentID] = append(newHierarchy[r.ParentID], r.ChildID)
	}

	PermissionHierarchy = newHierarchy
	slog.Info("Permission hierarchy loaded from DB", "parents", len(newHierarchy))
}

// CheckAccess перевіряє доступ, використовуючи регулярні вирази на рівні БД та кешування результатів у Redis.
func CheckAccess(method, path string) ([]string, bool, error) {
	cacheKey := fmt.Sprintf("auth:policy:%s:%s", strings.ToUpper(method), path)

	// 1. Спробуємо отримати з Redis
	if val, _ := RDB.Get(context.Background(), cacheKey).Result(); val != "" {
		if val == "UNMANAGED" {
			slog.Debug("Policy cache hit: unmanaged path", "path", path)
			return nil, false, nil // Не керується політиками
		}
		var reqPerms []string
		json.Unmarshal([]byte(val), &reqPerms)
		slog.Debug("Policy cache hit: rules found", "path", path, "required_perms", reqPerms)
		return reqPerms, true, nil // Керується, повертаємо список необхідних прав
	}

	slog.Info("Policy cache miss: querying database", "method", method, "path", path)

	// 2. Якщо в кеші немає — йдемо в Postgres
	// Спочатку перевіряємо, чи шлях взагалі керується хоч якимось правилом
	var allRulesForPath []models.PolicyRule
	if err := DB.Order("id asc").Where("? ~ path_pattern", path).Find(&allRulesForPath).Error; err != nil {
		return nil, false, err
	}

	if len(allRulesForPath) == 0 {
		// Шлях не керується ніякими правилами
		RDB.Set(context.Background(), cacheKey, "UNMANAGED", 15*time.Minute)
		return nil, false, nil
	}

	// 3. Фільтруємо правила за методом
	var requiredPerms []string
	for _, rule := range allRulesForPath {
		if rule.Method == "*" || strings.EqualFold(rule.Method, method) {
			for _, p := range strings.Split(rule.RequiredPermission, ",") {
				requiredPerms = append(requiredPerms, strings.TrimSpace(p))
			}
		} else if strings.ToUpper(rule.Method) == "CRUD" {
			action := "read"
			switch strings.ToUpper(method) {
			case "POST":
				action = "create"
			case "PUT", "PATCH":
				action = "update"
			case "DELETE":
				action = "delete"
			case "GET":
				action = "read"
			}
			for _, res := range strings.Split(rule.RequiredPermission, ",") {
				requiredPerms = append(requiredPerms, action+":"+strings.TrimSpace(res))
			}
		}
	}

	// 4. Кешуємо результат (навіть якщо список порожній — це означає "керований шлях, але метод не дозволений")
	val, _ := json.Marshal(requiredPerms)
	RDB.Set(context.Background(), cacheKey, string(val), 15*time.Minute)

	return requiredPerms, true, nil
}

// EvaluateAccess об’єднує перевірку політик та прав користувача
func EvaluateAccess(method, path string, userPerms []string) (bool, error) {
	reqPerms, isManaged, err := CheckAccess(method, path)
	if err != nil {
		return false, err
	}

	if !isManaged {
		return true, nil // Allow by default for unmanaged paths
	}

	// Якщо шлях керований, але жодне правило не підійшло під метод
	if len(reqPerms) == 0 {
		slog.Warn("Access denied: path is managed but method not allowed", "method", method, "path", path)
		return false, nil
	}
	// Перевіряємо, чи є у користувача хоча б одне з необхідних прав
	for _, rp := range reqPerms {
		if rp == "" || HasPermission(userPerms, rp) {
			return true, nil
		}
	}

	slog.Warn("Access denied: insufficient permissions", "method", method, "path", path, "required", reqPerms)
	return false, nil
}

// ValidatePermissions checks if all targetPerms are covered by userPerms (respecting hierarchy).
// This prevents permission escalation when assigning roles or keys.
func ValidatePermissions(userPerms []string, targetPerms []string) error {
	for _, tp := range targetPerms {
		if !HasPermission(userPerms, tp) {
			return fmt.Errorf("insufficient permission to grant: %s", tp)
		}
	}
	return nil
}

func HasPermission(userPerms []string, required string) bool {
	if required == "" {
		return true
	}

	for _, p := range userPerms {
		if p == "admin" || p == required {
			return true
		}

		// Перевірка через ієрархію (рекурсивно)
		if deps, ok := PermissionHierarchy[p]; ok {
			for _, dep := range deps {
				if HasPermission([]string{dep}, required) {
					return true
				}
			}
		}

		// Формат: action:resource[:scope] (наприклад, read:cameras:all або update:events)
		partsUser := strings.Split(p, ":")
		partsReq := strings.Split(required, ":")

		// Мінімально має бути action:resource
		if len(partsUser) < 2 || len(partsReq) < 2 {
			continue
		}

		actionUser := partsUser[0]
		resourceUser := partsUser[1]
		scopeUser := ""
		if len(partsUser) > 2 {
			scopeUser = partsUser[2]
		}

		actionReq := partsReq[0]
		resourceReq := partsReq[1]
		scopeReq := ""
		if len(partsReq) > 2 {
			scopeReq = partsReq[2]
		}

		// 1. Перевірка ресурсу (підтримка wildcard)
		if resourceUser != "*" && resourceUser != resourceReq {
			continue
		}

		// 2. Перевірка скоупу
		// Якщо у користувача 'all', він може все в межах ресурсу.
		// Якщо скоупи збігаються - ок.
		// Якщо у користувача немає 'all', а запитується конкретний скоуп - відмова (тільки якщо вони не однакові).
		scopeMatch := false
		if scopeUser == "all" {
			scopeMatch = true
		} else if scopeUser == scopeReq {
			scopeMatch = true
		}

		if !scopeMatch {
			continue
		}

		// 3. Ієрархія дій:
		// manage > delete > update/validate > create > read
		if actionUser == "manage" || actionUser == "admin" {
			return true
		}

		switch actionReq {
		case "read":
			// Будь-яка дія (create/update/delete/validate) дозволяє read
			return true
		case "create":
			if actionUser == "create" || actionUser == "update" || actionUser == "delete" || actionUser == "manage" {
				return true
			}
		case "update", "validate":
			// update та validate ми вважаємо на одному рівні або update сильніший
			if actionUser == "update" || actionUser == "validate" || actionUser == "delete" || actionUser == "manage" {
				return true
			}
		case "delete":
			if actionUser == "delete" || actionUser == "manage" {
				return true
			}
		}
	}
	return false
}

// ExpandPermissions розширює список прав користувача, додаючи всі залежні права згідно з ієрархією.
func ExpandPermissions(userPerms []string) []string {
	result := make(map[string]bool)
	var queue []string

	for _, p := range userPerms {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if !result[p] {
			result[p] = true
			queue = append(queue, p)
		}
	}

	for i := 0; i < len(queue); i++ {
		p := queue[i]

		// Додаємо базове право для прав з суфіксом :all (наприклад, manage:users:all -> manage:users)
		if strings.HasSuffix(p, ":all") {
			basePerm := strings.TrimSuffix(p, ":all")
			if !result[basePerm] {
				result[basePerm] = true
				queue = append(queue, basePerm)
			}

			// Якщо база має залежності, то ці залежності також можуть мати версію :all
			// Наприклад, manage:permits:all -> manage:permits -> read:permits
			// Отже ми маємо додати і read:permits:all
			if deps, ok := PermissionHierarchy[basePerm]; ok {
				for _, dep := range deps {
					depAll := dep + ":all"
					// Перевіряємо чи існує таке право взагалі (тільки якщо воно є в ієрархії або як база)
					// Але в ExpandPermissions ми зазвичай додаємо все що виглядає логічно
					if !result[depAll] {
						result[depAll] = true
						queue = append(queue, depAll)
					}
				}
			}
		}

		if deps, ok := PermissionHierarchy[p]; ok {
			for _, dep := range deps {
				if !result[dep] {
					result[dep] = true
					queue = append(queue, dep)
				}
			}
		}
	}

	final := make([]string, 0, len(result))
	for p := range result {
		final = append(final, p)
	}
	return final
}
