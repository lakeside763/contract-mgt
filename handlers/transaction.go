package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lakeside763/contract-mgt/models"
	"gorm.io/gorm"
)

func GetContracts(c *gin.Context, db *gorm.DB) {
	var contracts []models.Contract

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No user ID in context"})
		return
	}

	userRole, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No user role in context"})
		return
	}

	switch userRole {
	case "CONTRACTOR":
		query := `SELECT * FROM contracts WHERE "contractorId" = ?`
		if err := db.Where(query, userID).Find(&contracts).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contractor ID"})
			return
		}
	case "CLIENT":
		query := `SELECT * FROM contracts WHERE "clientId" = ?`
		if err := db.Raw(query, userID).Find(&contracts).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
			return
		}
	default:
		if err := db.Find(&contracts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, contracts)
}


func GetContract(c *gin.Context, db *gorm.DB) {
	contractID := c.Param("id")
	contract := models.Contract{ID: contractID}

	if err := db.First(&contract).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, contract)
}

func CreateContract(c *gin.Context, db *gorm.DB) {
	var contract models.Contract

	c.JSON(http.StatusOK, contract)
}