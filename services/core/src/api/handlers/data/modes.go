package data

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/truckguard/core/src/models"
	datarepo "github.com/truckguard/core/src/repository/data"
	"github.com/truckguard/core/src/utils"
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

func HandleUpdateMode(c *gin.Context) {
	id := c.Param("id")
	var input models.CustomsMode
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mode, err := datarepo.UpdateMode(c.Request.Context(), id, &input)
	if err != nil {
		slog.Error("Failed to update mode", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update mode: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, mode)
}

func HandleDeleteMode(c *gin.Context) {
	id := c.Param("id")
	if err := datarepo.DeleteMode(c.Request.Context(), id); err != nil {
		slog.Error("Failed to delete mode", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete mode: " + err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
