package main

import (
	"log"
	"os"

	"github.com/amarjeetdev/ev-charging-app/db"
	"github.com/amarjeetdev/ev-charging-app/models"
	"github.com/amarjeetdev/ev-charging-app/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Println("Error loading .env file")
	}

	r := gin.Default()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db.ConnectDatabase()
	err := db.DB.AutoMigrate(&models.Booking{}, &models.Station{}, &models.User{})

	if err != nil {
		log.Fatal("failed to migrate", err)
	}

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "hello world",
		})
	})

	routes.RegisterRoutes(r)

	log.Println("🤝🫱Server running on port " + port)
	r.Run(":" + port)
}
