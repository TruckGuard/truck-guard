package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/truckguard/customs-parser/src/api/handlers"
	"github.com/truckguard/customs-parser/src/api/middleware"
	"github.com/truckguard/customs-parser/src/pkg/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	logger := telemetry.NewLogger("truckguard-customs-parser")
	slog.SetDefault(logger)

	if err := telemetry.Init("truckguard-customs-parser"); err != nil {
		logger.Error("otel init failed", "error", err)
		os.Exit(1)
	}
	defer telemetry.Shutdown(context.Background())

	router := gin.New()
	router.Use(middleware.Logger(), gin.Recovery())
	router.Use(otelgin.Middleware("truckguard-customs-parser"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	customs := router.Group("/customs")
	{
		customs.GET("/declaration/:number", handlers.GetDeclarationMock)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	slog.Info("Customs Parser Service starting", "port", port)
	if err := router.Run(":" + port); err != nil {
		slog.Error("Failed to run server", "error", err)
		os.Exit(1)
	}
}
