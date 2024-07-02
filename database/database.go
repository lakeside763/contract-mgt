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
	Rdb *redis.Client
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

	// Initialize redis
	Rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Auto migrate User model
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.Contract{},
		&models.Job{},
	); err != nil {
		return err
	}

	return nil
}