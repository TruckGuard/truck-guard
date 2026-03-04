package data

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/truckguard/core/src/models"
	datarepo "github.com/truckguard/core/src/repository/data"
	"github.com/truckguard/core/src/utils"
)

func HandleListVehicleTypes(c *gin.Context) {
	limit, offset, page := utils.GetPagination(c)
	filter := datarepo.VehicleTypeFilter{
		Code: c.Query("code"),
		Name: c.Query("name"),
	}
	types, total, err := datarepo.ListVehicleTypes(c.Request.Context(), limit, offset, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vehicle types"})
		return
	}
	utils.SendPaginatedResponse(c, types, total, page, limit)
}

func HandleCreateVehicleType(c *gin.Context) {
	var input models.VehicleType
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := datarepo.CreateVehicleType(c.Request.Context(), &input); err != nil {
		slog.Error("Failed to create vehicle type", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vehicle type: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, input)
}

func HandleUpdateVehicleType(c *gin.Context) {
	id := c.Param("id")
	var input models.VehicleType
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vt, err := datarepo.UpdateVehicleType(c.Request.Context(), id, &input)
	if err != nil {
		slog.Error("Failed to update vehicle type", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vehicle type: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, vt)
}

func HandleDeleteVehicleType(c *gin.Context) {
	id := c.Param("id")
	if err := datarepo.DeleteVehicleType(c.Request.Context(), id); err != nil {
		slog.Error("Failed to delete vehicle type", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vehicle type: " + err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
