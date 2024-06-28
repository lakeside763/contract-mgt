package database

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/lakeside763/contract-mgt/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	DB *gorm.DB

	Rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
)

func InitDB() error {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy {
			NameReplacer: strings.NewReplacer("ID", "id"),
		},
	})
	if err != nil {
		return err
	}

	// Auto migrate User model
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Profile{})
	DB.AutoMigrate(models.Contract{})
	DB.AutoMigrate(models.Job{})
	return nil
}