package data

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/truckguard/core/src/models"
	"github.com/truckguard/core/src/repository"
	"github.com/truckguard/core/src/utils"
	"gorm.io/gorm"
)

func HandleListExcludedPlates(c *gin.Context) {
	limit, offset, page := utils.GetPagination(c)
	var plates []models.ExcludedPlate
	var total int64

	db := repository.DB.WithContext(c.Request.Context())
	
	if err := db.Model(&models.ExcludedPlate{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка при підрахунку ігнорованих номерів"})
		return
	}

	if err := db.Limit(limit).Offset(offset).Find(&plates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка при отриманні списку ігнорованих номерів"})
		return
	}

	utils.SendPaginatedResponse(c, plates, total, page, limit)
}

func HandleGetExcludedPlateByID(c *gin.Context) {
	id := c.Param("id")
	var plate models.ExcludedPlate
	if err := repository.DB.WithContext(c.Request.Context()).First(&plate, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ігнорований номер не знайдено"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка при отриманні ігнорованого номера"})
		}
		return
	}
	c.JSON(http.StatusOK, plate)
}

func HandleCreateExcludedPlate(c *gin.Context) {
	var input models.ExcludedPlate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := repository.DB.WithContext(c.Request.Context()).Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка при додаванні номера до списку ігнорування"})
		return
	}

	c.JSON(http.StatusCreated, input)
}

func HandleUpdateExcludedPlate(c *gin.Context) {
	id := c.Param("id")
	var plate models.ExcludedPlate
	db := repository.DB.WithContext(c.Request.Context())

	if err := db.First(&plate, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ігнорований номер не знайдено"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка при пошуку ігнорованого номера"})
		}
		return
	}

	if err := c.ShouldBindJSON(&plate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Save(&plate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка при оновленні ігнорованого номера"})
		return
	}

	c.JSON(http.StatusOK, plate)
}

func HandleDeleteExcludedPlate(c *gin.Context) {
	id := c.Param("id")
	db := repository.DB.WithContext(c.Request.Context())

	var plate models.ExcludedPlate
	if err := db.First(&plate, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ігнорований номер не знайдено"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка при пошуку ігнорованого номера"})
		}
		return
	}

	if err := db.Delete(&plate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка при видаленні номера зі списку ігнорування"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
