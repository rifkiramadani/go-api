package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	// Route
	server.GET("/events", getEvents)
	server.GET("/test", getTest)

	server.Run(":8080")
}

// Function Handler
func getEvents(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "Hello From Event",
	})
}

func getTest(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "Testing",
	})
}
