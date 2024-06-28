package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/lakeside763/contract-mgt/handlers"
	"github.com/lakeside763/contract-mgt/middlewares"
	"gorm.io/gorm"
)

func ProfileRoutes(router *gin.Engine, db *gorm.DB, authorized *gin.RouterGroup, rdb *redis.Client) {
	authorized.Use(middlewares.AuthMiddleware(rdb)) 
	{
		authorized.GET("/profiles", func(c *gin.Context) {handlers.GetProfiles(c, db)})
		authorized.GET("/profiles/:id", func(c *gin.Context) {handlers.GetProfile(c, db)})
		authorized.POST("/profiles", func(c *gin.Context) {handlers.CreateProfile(c, db)})
	}
}