package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lakeside763/contract-mgt/database"
	"github.com/lakeside763/contract-mgt/routes"
)



func main() {
	if err := database.InitDB(); err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	rdb := database.Rdb
	router := gin.Default()
	

	// routes
	routes.UserRoutes(router, database.DB)
	routes.AuthRoutes(router, database.DB, rdb)
	authorized := router.Group("/")
	routes.ProfileRoutes(router, database.DB, authorized, rdb)
	routes.TransactionRoutes(router, database.DB, authorized, rdb)

	router.Run(":5200")
}