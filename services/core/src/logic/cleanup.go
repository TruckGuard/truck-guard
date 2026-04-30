package logic

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/truckguard/core/src/repository"
)

// RunCleanupBackground starts a background goroutine that periodically cleans up the database.
func RunCleanupBackground(ctx context.Context) {
	slog.Info("Starting background cleanup scheduler...")

	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		// Run immediately on start
		performCleanup(ctx)

		for {
			select {
			case <-ctx.Done():
				slog.Info("Stopping background cleanup scheduler")
				ticker.Stop()
				return
			case <-ticker.C:
				performCleanup(ctx)
			}
		}
	}()
}

func performCleanup(ctx context.Context) {
	slog.Info("Running periodic database cleanup...")

	// 1. Void Unverified Permits
	hoursStr := repository.GetSystemSetting(ctx, "unverified_permits_cleanup_hours")
	if hoursStr != "" && hoursStr != "0" {
		hours, err := strconv.Atoi(hoursStr)
		if err == nil && hours > 0 {
			rows, err := repository.VoidUnverifiedPermits(ctx, hours)
			if err != nil {
				slog.Error("Failed to void unverified permits", "error", err)
			} else if rows > 0 {
				slog.Info("Voided unverified permits", "count", rows)
			}
		}
	}

	// 2. Cleanup Audit Logs
	daysAuditStr := repository.GetSystemSetting(ctx, "permit_audit_cleanup_days")
	if daysAuditStr != "" && daysAuditStr != "0" {
		days, err := strconv.Atoi(daysAuditStr)
		if err == nil && days > 0 {
			rows, err := repository.CleanupAuditLogs(ctx, days)
			if err != nil {
				slog.Error("Failed to cleanup audit logs", "error", err)
			} else if rows > 0 {
				slog.Info("Cleaned up audit logs", "count", rows)
			}
		}
	}

	// 3. Cleanup Voided Permits
	daysVoidStr := repository.GetSystemSetting(ctx, "void_permits_deletion_days")
	if daysVoidStr != "" && daysVoidStr != "0" {
		days, err := strconv.Atoi(daysVoidStr)
		if err == nil && days > 0 {
			rows, err := repository.CleanupVoidedPermits(ctx, days)
			if err != nil {
				slog.Error("Failed to cleanup voided permits", "error", err)
			} else if rows > 0 {
				slog.Info("Permanently deleted old voided permits", "count", rows)
			}
		}
	}

	// 4. Cleanup Notifications
	daysNotifStr := repository.GetSystemSetting(ctx, "notification_retention_days")
	if daysNotifStr != "" && daysNotifStr != "0" {
		days, err := strconv.Atoi(daysNotifStr)
		if err == nil && days > 0 {
			rows, err := repository.CleanupNotifications(ctx, days)
			if err != nil {
				slog.Error("Failed to cleanup notifications", "error", err)
			} else if rows > 0 {
				slog.Info("Permanently deleted old notifications", "count", rows)
			}
		}
	}

	slog.Info("Periodic database cleanup completed")
}
