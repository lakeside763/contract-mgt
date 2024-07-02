package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lakeside763/contract-mgt/handlers"
	"gorm.io/gorm"
)

// func CreateUser(router *gin.Engine) {
// 	router.POST("/users", handlers.CreateUser)
// }


func UserRoutes(router *gin.Engine, db *gorm.DB ) {
	users := router.Group("/v1/api/users")
	{
		users.GET("/", func(c *gin.Context) { handlers.GetUsers(c, db) }) 
		users.POST("/", func(c *gin.Context) { handlers.CreateUser(c, db) })
		users.GET("/:id", func(c *gin.Context) {handlers.GetUser(c, db) })
	}
}