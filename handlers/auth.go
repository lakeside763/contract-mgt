package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lakeside763/contract-mgt/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)


func Login(c *gin.Context, db *gorm.DB, rdb *redis.Client) {
	var creds models.Credentials

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.Where("username = ?", creds.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}

	expirationTime := time.Now().Add(14 * 24 * time.Hour)
	claims := &models.Claims {
		ID: user.ID,
		Username: user.Username,
		Role: string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims {
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(models.JwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	if _, err := verifyToken(tokenString); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token generated"})
		return
	}

	err = rdb.Set(c, tokenString, "valid", expirationTime.Sub(time.Now().Add(14 * 24 * time.Hour))).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not store the token"})
	}

	userWithToken := models.UserResponseWithToken {
		UserResponse: models.UserResponse {
			ID:        user.ID,
			Username:  user.Username,
			Role:      string(user.Role),
			CreatedAt: user.CreatedAt,
		},
		Token: tokenString,
	}

	c.JSON(http.StatusOK, userWithToken)
}

func Logout(c *gin.Context, db *gorm.DB, rdb *redis.Client) {
	authHeader := c.GetHeader("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
		return
	}

	tokenString := parts[1]

	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
		return
	}

	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return models.JwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
	}

	// claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-1 * time.Hour))

	// _, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(models.JwtKey)
	// if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not expire token"})
	// 		return
	// }

	err = rdb.Del(c, tokenString).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not remove from redis"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func verifyToken(tokenString string) (*models.Claims, error) {
	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return models.JwtKey, nil
	})
	
	// Log verification process
	if err != nil {
		log.Printf("Token verification failed: %v", err)
		return nil, err
	}

	if !token.Valid {
		log.Println("Token is invalid")
		return nil, fmt.Errorf("token is invalid")
	}

	log.Println("Token is valid")
	
	return claims, nil
}