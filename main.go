package main

import (
	"ecommerce-backend/db"
	"ecommerce-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	defer db.DB.Close()

	// Initialize the Gin router
	router := gin.Default()

	// Configure CORS middleware
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:5173"} // Allow requests from SvelteKit dev server
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	// Setup routes
	routes.SetupRoutes(router)

	// Start the server
	router.Run(":8080")
}
