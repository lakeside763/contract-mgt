package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/lakeside763/contract-mgt/handlers"
	"github.com/lakeside763/contract-mgt/middlewares"
	"gorm.io/gorm"
)


func TransactionRoutes(router *gin.Engine, db *gorm.DB, rdb *redis.Client) {
	// Group routes under /v1/api/transactions
	transactions := router.Group("/v1/api/transactions")
	{
		// use the authorization middleware for the group
		transactions.Use(middlewares.AuthMiddleware(rdb))

		// Define contract-related routes
		transactions.GET("/contracts", func(c *gin.Context) {handlers.GetContracts(c, db)})
		transactions.GET("/contracts/:id", func(c *gin.Context) {handlers.GetContract(c, db)})
		transactions.POST("/contracts", func(c *gin.Context) {handlers.CreateContract(c, db)})

		// Define job-related routes
		transactions.POST("/jobs", func(c *gin.Context) {handlers.CreateJob(c, db)})
		transactions.GET("/jobs", func(c *gin.Context) {handlers.GetJobs(c, db)})

		// Define payment-related routes
		transactions.POST("/jobs/payment", func(c *gin.Context) {handlers.JobPayment(c, db)})
		transactions.POST("/payments/deposit", func(c *gin.Context) {handlers.PaymentDeposit(c, db)})
	}
}

