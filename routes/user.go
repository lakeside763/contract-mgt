package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lakeside763/contract-mgt/handlers"
	"gorm.io/gorm"
)

// func CreateUser(router *gin.Engine) {
// 	router.POST("/users", handlers.CreateUser)
// }


func UserRoutes(router *gin.Engine, db *gorm.DB) {
	router.GET("/users", func(c *gin.Context) { handlers.GetUsers(c, db) }) 
	router.POST("/users", func(c *gin.Context) { handlers.CreateUser(c, db) })
	router.GET("/users/:id", func(c *gin.Context) {handlers.GetUser(c, db) })
}