package main

import (
	"log"
	"net/http"

	"example.com/event-app/config"
	"example.com/event-app/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()

	server := gin.Default()

	// Route
	api := server.Group("/api")
	{
		api.POST("/events", createEvent)
		api.GET("/events", getEvents)
	}

	server.Run(":8080")
}

// Function Handler
func getEvents(context *gin.Context) {
	events := models.GetAllEvents()

	context.JSON(http.StatusOK, events)
}

func createEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request data",
			"error":   err.Error(),
		})
		return
	}

	// Dummy Data
	event.UserId = 1

	// Save Inputan
	event.Save()

	context.JSON(http.StatusCreated, gin.H{
		"message": "Create Event",
		"event":   event,
	})
}
