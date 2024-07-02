package services

import (
	"github.com/lakeside763/contract-mgt/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateProfileAndUser(db *gorm.DB, profile models.Profile, creds models.Credentials) (models.Profile, error) {
	err := db.Transaction(func(tx *gorm.DB) error {
		// Create Profile
		if err := tx.Create(&profile).Error; err != nil {
			return err
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		// Create User
		user := models.User{
			ID:       profile.ID,
			Username: creds.Username,
			Password: string(hashedPassword),
			Role:     models.Role(profile.Type),
		}

		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		return nil
	})

	return profile, err
}