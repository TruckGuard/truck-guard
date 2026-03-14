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

func HandleListPaymentTypes(c *gin.Context) {
	limit, offset, page := utils.GetPagination(c)
	filter := datarepo.PaymentTypeFilter{
		Code:     c.Query("code"),
		IsActive: c.Query("is_active"),
	}
	types, total, err := datarepo.ListPaymentTypes(c.Request.Context(), limit, offset, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payment types"})
		return
	}
	utils.SendPaginatedResponse(c, types, total, page, limit)
}

func HandleCreatePaymentType(c *gin.Context) {
	var input models.PaymentType
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := datarepo.CreatePaymentType(c.Request.Context(), &input); err != nil {
		slog.Error("Failed to create payment type", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment type: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, input)
}

func HandleGetPaymentTypeByID(c *gin.Context) {
	id := c.Param("id")
	paymentType, err := datarepo.GetPaymentType(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Тип оплати не знайдено"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка при отриманні типу оплати"})
		}
		return
	}
	c.JSON(http.StatusOK, paymentType)
}

func HandleUpdatePaymentType(c *gin.Context) {
	id := c.Param("id")
	var input models.PaymentType
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pt, err := datarepo.UpdatePaymentType(c.Request.Context(), id, &input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Тип оплати не знайдено"})
		} else {
			slog.Error("Failed to update payment type", "error", err, "id", id)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка при оновленні типу оплати: " + err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, pt)
}

func HandleDeletePaymentType(c *gin.Context) {
	id := c.Param("id")
	if err := datarepo.DeletePaymentType(c.Request.Context(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Тип оплати не знайдено"})
		} else {
			slog.Error("Failed to delete payment type", "error", err, "id", id)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка при видаленні типу оплати: " + err.Error()})
		}
		return
	}
	c.Status(http.StatusNoContent)
}
