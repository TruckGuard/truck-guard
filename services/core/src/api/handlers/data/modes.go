package data

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/truckguard/core/src/models"
	datarepo "github.com/truckguard/core/src/repository/data"
	"github.com/truckguard/core/src/utils"
	"gorm.io/gorm"
)

func HandleListModes(c *gin.Context) {
	limit, offset, page := utils.GetPagination(c)
	filter := datarepo.ModeFilter{
		Code: c.Query("code"),
		Name: c.Query("name"),
	}
	modes, total, err := datarepo.ListModes(c.Request.Context(), limit, offset, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch modes"})
		return
	}
	utils.SendPaginatedResponse(c, modes, total, page, limit)
}

func HandleCreateMode(c *gin.Context) {
	var input models.CustomsMode
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := datarepo.CreateMode(c.Request.Context(), &input); err != nil {
		slog.Error("Failed to create mode", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create mode: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, input)
}

func HandleGetModeByID(c *gin.Context) {
	id := c.Param("id")
	mode, err := datarepo.GetMode(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Митний режим не знайдено"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка при отриманні митного режиму"})
		}
		return
	}
	c.JSON(http.StatusOK, mode)
}

func HandleUpdateMode(c *gin.Context) {
	id := c.Param("id")
	var input models.CustomsMode
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mode, err := datarepo.UpdateMode(c.Request.Context(), id, &input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Митний режим не знайдено"})
		} else {
			slog.Error("Failed to update mode", "error", err, "id", id)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка при оновленні митного режиму: " + err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, mode)
}

func HandleDeleteMode(c *gin.Context) {
	id := c.Param("id")
	if err := datarepo.DeleteMode(c.Request.Context(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Митний режим не знайдено"})
		} else {
			slog.Error("Failed to delete mode", "error", err, "id", id)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка при видаленні митного режиму: " + err.Error()})
		}
		return
	}
	c.Status(http.StatusNoContent)
}
