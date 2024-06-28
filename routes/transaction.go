package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/lakeside763/contract-mgt/handlers"
	"github.com/lakeside763/contract-mgt/middlewares"
	"gorm.io/gorm"
)


func TransactionRoutes(router *gin.Engine, db *gorm.DB, authorized *gin.RouterGroup, rdb *redis.Client) {
	authorized.Use(middlewares.AuthMiddleware(rdb))
	{
		authorized.GET("/transactions/contracts", func(c *gin.Context) {handlers.GetContracts(c, db)})
		authorized.GET("/transactions/contracts/:id", func(c *gin.Context) {handlers.GetContract(c, db)})
	}
}

