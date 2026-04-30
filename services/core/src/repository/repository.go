package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/truckguard/core/src/models"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dsn string) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to Core Database")
	}

	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		panic(err)
	}

	// Drop old strict constraints to allow GORM to rebuild them with CASCADE
	db.Exec(`ALTER TABLE "plate_events" DROP CONSTRAINT IF EXISTS "fk_plate_events_camera"`)
	db.Exec(`ALTER TABLE "weight_events" DROP CONSTRAINT IF EXISTS "fk_weight_events_scale"`)

	db.AutoMigrate(
		&models.CustomsPost{},
		&models.SystemEvent{},
		&models.PlateEvent{},
		&models.WeightEvent{},
		&models.CameraConfig{},
		&models.ScaleConfig{},
		&models.SystemSetting{},
		&models.ExcludedPlate{},
		&models.Permit{},
		&models.PermitAudit{},
		&models.User{},
		&models.CustomsMode{},
		&models.Company{},
		&models.VehicleType{},
		&models.PaymentType{},
		&models.PermitPayer{},
		&models.Notification{},
		&models.SystemAudit{},
	)
	DB = db
}

var RDB *redis.Client

func InitRedis(addr string) {
	RDB = redis.NewClient(&redis.Options{
		Addr:         addr,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})
	if _, err := RDB.Ping(context.Background()).Result(); err != nil {
		panic("Failed to connect to Redis")
	}
}
