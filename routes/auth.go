package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/lakeside763/contract-mgt/handlers"
	"gorm.io/gorm"
)

func AuthRoutes(router *gin.Engine, db *gorm.DB, rdb *redis.Client) {
	auth := router.Group("/v1/api") 
	{
		auth.POST("/login", func(c *gin.Context) { handlers.Login(c, db, rdb)})
		auth.POST("/logout", func(c *gin.Context) { handlers.Logout(c, db,rdb)})
	}
}