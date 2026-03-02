package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/truckguard/core/src/api/clients"
	"github.com/truckguard/core/src/models"
	"github.com/truckguard/core/src/repository"
	"github.com/truckguard/core/src/utils"
)

func HandleCreateCamera(c *gin.Context) {
	var config models.CameraConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authClient := clients.NewAuthClient(c.GetHeader("Authorization"), c.GetHeader("X-Api-Key"))
	authResp, err := authClient.CreateApiKey(
		c.Request.Context(),
		config.Name,
		[]string{"ingest:events"},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create authentication key: " + err.Error()})
		return
	}

	if config.SourceID == "" {
		config.SourceID = fmt.Sprintf("%v", authResp.ID)
	}

	slog.Debug("Camera config created", "config", config, "source_id", config.SourceID)

	if err := repository.CreateCamera(c.Request.Context(), &config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save camera configuration"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"camera":  config,
		"api_key": authResp.APIKey,
	})
}

func HandleGetCameras(c *gin.Context) {
	limit, offset, page := utils.GetPagination(c)

	configs, total, err := repository.GetCameras(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cameras"})
		return
	}

	utils.SendPaginatedResponse(c, configs, total, page, limit)
}

func HandleGetConfigByID(c *gin.Context) {
	sourceID := c.Param("id")
	config, err := repository.GetCameraByID(c.Request.Context(), sourceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Camera config not found"})
		return
	}
	c.JSON(http.StatusOK, config)
}

func HandleGetConfigByCameraID(c *gin.Context) {
	sourceID := c.Param("camera_id")
	config, err := repository.GetCameraBySourceID(c.Request.Context(), sourceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Camera config not found"})
		return
	}
	c.JSON(http.StatusOK, config)
}

func HandleUpdateCamera(c *gin.Context) {
	id := c.Param("id")
	config, err := repository.GetCameraByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Camera configuration not found"})
		return
	}

	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := repository.UpdateCamera(c.Request.Context(), &config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update configuration"})
		return
	}

	c.JSON(http.StatusOK, config)
}

func HandleDeleteCamera(c *gin.Context) {
	id := c.Param("id")
	config, err := repository.GetCameraByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Camera config not found"})
		return
	}

	if config.SourceID != "" {
		authClient := clients.NewAuthClient(c.GetHeader("Authorization"), c.GetHeader("X-Api-Key"))
		err := authClient.DeleteApiKey(
			c.Request.Context(),
			config.SourceID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete associated API key: " + err.Error()})
			return
		}
	}

	if err := repository.DeleteCamera(c.Request.Context(), &config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete camera configuration"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func HandleCreateScale(c *gin.Context) {
	var config models.ScaleConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authClient := clients.NewAuthClient(c.GetHeader("Authorization"), c.GetHeader("X-Api-Key"))
	authResp, err := authClient.CreateApiKey(
		c.Request.Context(),
		config.Name+"_key",
		[]string{"ingest:events"},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create authentication key: " + err.Error()})
		return
	}

	if config.SourceID == "" {
		config.SourceID = fmt.Sprintf("%v", authResp.ID)
	}

	slog.Debug("Scale config created", "config", config, "match_permit", config.MatchPermit)
	if err := repository.CreateScale(c.Request.Context(), &config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save scale configuration"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"scale":   config,
		"api_key": authResp.APIKey,
	})
}

func HandleGetScales(c *gin.Context) {
	limit, offset, page := utils.GetPagination(c)

	configs, total, err := repository.GetScales(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch scales"})
		return
	}

	utils.SendPaginatedResponse(c, configs, total, page, limit)
}

func HandleGetConfigByScaleID(c *gin.Context) {
	scaleID := c.Param("scale_id")
	config, err := repository.GetScaleBySourceID(c.Request.Context(), scaleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "scale config not found"})
		return
	}
	c.JSON(http.StatusOK, config)
}

func HandleUpdateScale(c *gin.Context) {
	id := c.Param("id")
	config, err := repository.GetScaleByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Scale configuration not found"})
		return
	}

	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := repository.UpdateScale(c.Request.Context(), &config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update scale configuration"})
		return
	}

	c.JSON(http.StatusOK, config)
}

func HandleDeleteScale(c *gin.Context) {
	id := c.Param("id")
	config, err := repository.GetScaleByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Scale configuration not found"})
		return
	}

	if config.SourceID != "" {
		authClient := clients.NewAuthClient(c.GetHeader("Authorization"), c.GetHeader("X-Api-Key"))
		err := authClient.DeleteApiKey(
			c.Request.Context(),
			config.SourceID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete associated API key: " + err.Error()})
			return
		}
	}

	if err := repository.DeleteScale(c.Request.Context(), &config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete scale configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
 