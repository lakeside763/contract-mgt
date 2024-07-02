package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lakeside763/contract-mgt/database"
	"github.com/lakeside763/contract-mgt/models"
	"github.com/lakeside763/contract-mgt/routes"
)



func main() {
	if err := database.InitDB(); err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	models.Init()

	rdb := database.Rdb
	router := gin.Default()
	

	// routes
	routes.UserRoutes(router, database.DB)
	routes.AuthRoutes(router, database.DB, rdb)
	routes.ProfileRoutes(router, database.DB, rdb)
	routes.TransactionRoutes(router, database.DB, rdb)

	
	// start the server with graceful shutdown
	runServer(router)
	// router.Run(":5200")
}

func runServer(router *gin.Engine) {
	const port = ":5200"
	server := &http.Server{
		Addr:    port,
		Handler: router,
	}

	// Channel to listen for termination signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v", err)
		}
	}()

	// Wait for a termination signal
	<-quit
	log.Println("Shutting down server...")

	// Create a context with timeout for the shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}