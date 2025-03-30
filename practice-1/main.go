package main

import (
	"log"
	"os"
	_ "practice-1/docs" // This is where the generated swagger docs will be
	"practice-1/soap"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Practice 1 API
// @version         1.0
// @description     A simple API using Gin framework with REST and SOAP endpoints
// @host           localhost:8080
// @BasePath       /api/v1

func init() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

// @Summary      Health check endpoint
// @Description  Returns the health status of the API
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /api/v1/health [get]
func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func initializeRoutes(router *gin.Engine) {
	// API v1 group
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", healthCheck)
	}

	// SOAP endpoints
	soapGroup := router.Group("/soap")
	{
		soapGroup.POST("/convert-currency", soap.HandleCurrencyConversion)
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func main() {
	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Initialize routes
	initializeRoutes(router)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
