package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/truckguard/ingestor/src/api/handlers"
	"github.com/truckguard/ingestor/src/api/middleware"
	"github.com/truckguard/ingestor/src/pkg/telemetry"
	"github.com/truckguard/ingestor/src/repository"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	logger := telemetry.NewLogger("truckguard-ingestor")
	slog.SetDefault(logger)

	if err := telemetry.Init("truckguard-ingestor"); err != nil {
		logger.Error("otel init failed", "error", err)
		os.Exit(1)
	}
	defer telemetry.Shutdown(context.Background())

	valkeyAddr := os.Getenv("VALKEY_ADDR")
	if valkeyAddr == "" {
		valkeyAddr = os.Getenv("REDIS_ADDR")
	}
	repository.InitRedis(valkeyAddr)

	endpoint := os.Getenv("STORAGE_ENDPOINT")
	accessKey := os.Getenv("STORAGE_ACCESS_KEY")
	secretKey := os.Getenv("STORAGE_SECRET_KEY")
	repository.InitMinio(endpoint, accessKey, secretKey)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(otelgin.Middleware("truckguard-ingestor"))

	ingestLines := r.Group("/ingest")
	{
		ingestLines.POST("/camera", handlers.HandleCameraIngest)
		ingestLines.POST("/weight", handlers.HandleWeightIngest)
	}

	r.Match([]string{"GET", "HEAD"}, "/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
