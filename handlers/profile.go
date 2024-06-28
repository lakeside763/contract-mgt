package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lakeside763/contract-mgt/models"
	"github.com/lakeside763/contract-mgt/services"
	"gorm.io/gorm"
)

func GetProfiles(c *gin.Context, db *gorm.DB) {
	var profiles []models.Profile

	if err := db.Find(&profiles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profiles)
}

func GetProfile(c *gin.Context, db *gorm.DB) {
	var profileID = c.Param("id")

	var profile = models.Profile{ID: profileID}
	if err := db.First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func CreateProfile(c *gin.Context, db *gorm.DB) {
	var puc models.ProfileAndUserCredentials
	
	if err := c.ShouldBindJSON(&puc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile := models.Profile{
		FirstName: puc.FirstName,
		LastName: puc.LastName,
		Profession: puc.Profession,
		Balance: puc.Balance,
		Type: puc.Type,
	}

	creds := models.Credentials{
		Username: puc.Username,
		Password: puc.Password,
	}

	profile, err := services.CreateProfileAndUser(db, profile, creds); 
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}