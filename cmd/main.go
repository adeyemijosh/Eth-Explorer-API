package main

import (
	"fmt"
	"log"
	"net/http"

	"eth-explorer-api/internal/config"
	"eth-explorer-api/internal/handlers"
	"eth-explorer-api/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("=== Starting Ethereum Explorer API ===")

	fmt.Println("Loading configuration...")
	cfg := config.Load()
	fmt.Printf("Config loaded - Port: %s, ETH_NODE_URL: %s\n", cfg.Port, cfg.EthNodeURL)

	fmt.Println("Initializing Ethereum service...")
	ethService, err := services.NewEthService(cfg.EthNodeURL)
	if err != nil {
		fmt.Printf("ERROR: Failed to initialize Ethereum service: %v\n", err)
		log.Fatal("Failed to initialize Ethereum service:", err)
	}
	fmt.Println("Ethereum service initialized successfully!")

	fmt.Println("Initializing handlers...")
	ethHandler := handlers.NewEthHandler(ethService)
	fmt.Println("Handlers initialized!")

	fmt.Println("Setting up Gin router...")
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	api := router.Group("/api/v1")
	{
		// Ethereum endpoints
		api.GET("/eth/block/:number", ethHandler.GetBlock)
		api.GET("/eth/transaction/:hash", ethHandler.GetTransaction)
		api.GET("/eth/balance/:address", ethHandler.GetBalance)
		api.GET("/eth/latest-block", ethHandler.GetLatestBlock)
		api.GET("/eth/gas-price", ethHandler.GetGasPrice)

		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "healthy"})
		})
	}

	fmt.Printf("Starting server on port %s...\n", cfg.Port)
	fmt.Println("=== Server should be running now ===")
	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(router.Run(":" + cfg.Port))
}
