package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/lakeside763/contract-mgt/handlers"
	"github.com/lakeside763/contract-mgt/middlewares"
	"gorm.io/gorm"
)

func ProfileRoutes(router *gin.Engine, db *gorm.DB, rdb *redis.Client) {
	profiles := router.Group("/v1/api/profiles")
	{
		profiles.Use(middlewares.AuthMiddleware(rdb)) 
		profiles.GET("/", func(c *gin.Context) {handlers.GetProfiles(c, db)})
		profiles.GET("/:id", func(c *gin.Context) {handlers.GetProfile(c, db)})
		profiles.POST("/", func(c *gin.Context) {handlers.CreateProfile(c, db)})
	}
}