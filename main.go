package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	router := gin.Default()

	// Initialize PostgreSQL connection
	_, sqlDB, err := initDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer sqlDB.Close()

	// Simple ping test to ensure service is running
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Route to hit the order service ping endpoint
	router.GET("/hit-order-ping", func(c *gin.Context) {
		resp, err := http.Get("http://edot-order-service:8080/ping")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to contact order service",
			})
			return
		}
		defer resp.Body.Close()

		// Forward the response from the order service
		if resp.StatusCode == http.StatusOK {
			c.JSON(http.StatusOK, gin.H{
				"message": "Successfully hit order service ping",
			})
		} else {
			c.JSON(resp.StatusCode, gin.H{
				"error": "Bad response from order service",
			})
		}
	})

	router.Run(":8080")
}

func initDB() (*gorm.DB, *sql.DB, error) {
	// PostgreSQL DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db, sqlDB, nil
}
