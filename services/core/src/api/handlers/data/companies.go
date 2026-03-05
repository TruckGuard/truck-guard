package data

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/truckguard/core/src/models"
	datarepo "github.com/truckguard/core/src/repository/data"
	"github.com/truckguard/core/src/utils"
)

func HandleListCompanies(c *gin.Context) {
	limit, offset, page := utils.GetPagination(c)
	filter := datarepo.CompanyFilter{
		Name:   c.Query("name"),
		EDRPOU: c.Query("edrpou"),
		Search: c.Query("q"),
	}
	companies, total, err := datarepo.ListCompanies(c.Request.Context(), limit, offset, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch companies"})
		return
	}
	utils.SendPaginatedResponse(c, companies, total, page, limit)
}

func HandleCreateCompany(c *gin.Context) {
	var input models.Company
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := datarepo.CreateCompany(c.Request.Context(), &input); err != nil {
		slog.Error("Failed to create company", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create company: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, input)
}

func HandleUpdateCompany(c *gin.Context) {
	id := c.Param("id")
	var input models.Company
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	company, err := datarepo.UpdateCompany(c.Request.Context(), id, &input)
	if err != nil {
		slog.Error("Failed to update company", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update company: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, company)
}

func HandleDeleteCompany(c *gin.Context) {
	id := c.Param("id")
	if err := datarepo.DeleteCompany(c.Request.Context(), id); err != nil {
		slog.Error("Failed to delete company", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete company: " + err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
