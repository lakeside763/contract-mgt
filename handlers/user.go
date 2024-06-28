package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lakeside763/contract-mgt/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(c *gin.Context, db *gorm.DB) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User
	err := db.Where("username = ?", user.Username).First(&existingUser).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user.Password = string(hashedPassword)
	user.Role = models.ADMIN

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := models.UserResponse {
		ID:					user.ID,
		Username: 	user.Username,
		Role: 			string(user.Role),
		CreatedAt: 	user.CreatedAt,
		UpdatedAt: 	user.UpdatedAt,	
	}
	
	c.JSON(http.StatusOK, response)
}

func GetUser(c *gin.Context, db *gorm.DB) {
	userID := c.Param("id")

	var user = models.User{ID: userID}

	if err := db.First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := models.UserResponse {
		ID:					user.ID,
		Username: 	user.Username,
		Role: 			string(user.Role),
		CreatedAt: 	user.CreatedAt,
		UpdatedAt: 	user.UpdatedAt,	
	}

	c.JSON(http.StatusOK, response)
}

func GetUsers(c *gin.Context, db *gorm.DB) {
	var users  []models.User
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responseUsers []models.UserResponse
	for _, user := range users {
		responseUser := models.UserResponse {
			ID:					user.ID,
			Username: 	user.Username,
			Role: 			string(user.Role),
			CreatedAt: 	user.CreatedAt,
			UpdatedAt: 	user.UpdatedAt,	
		}
		responseUsers = append(responseUsers, responseUser)
	}

	c.JSON(http.StatusOK, responseUsers)
}