package controllers

import (
	"net/http"

	"example.com/event-app/config"
	"example.com/event-app/models"
	"github.com/gin-gonic/gin"
)

func CreateEvent(context *gin.Context) {
	userID, _ := context.Get("userID")

	var event models.Event

	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//dummy userID
	event.UserID = uint(userID.(int))

	config.DB.Create(&event)
	context.JSON(http.StatusCreated, gin.H{
		"message": "Data Berhasil Di Buat",
		"event":   event,
	})
}

func GetEvents(context *gin.Context) {
	var events []models.Event

	config.DB.Find(&events)
	context.JSON(http.StatusOK, gin.H{
		"message": "Berhasil Menampilkan Semua Data Events",
		"events":  events,
	})
}

func GetEventbyId(context *gin.Context) {
	var event models.Event
	paramsId := context.Param("id")

	var eventData = config.DB.First(&event, paramsId).Error

	if eventData != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "Event Tidak Ditemukan",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"messsage": "Berhasil Menampilkan Detail Event",
		"event":    event,
	})
}

func UpdateEvent(context *gin.Context) {
	userID, _ := context.Get("userID")

	var event models.Event
	paramsId := context.Param("id")

	var eventData = config.DB.First(&event, paramsId).Error
	if eventData != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "Event Tidak Ditemukan",
		})
		return
	}

	if event.UserID != uint(userID.(int)) {
		context.JSON(http.StatusForbidden, gin.H{
			"error": "Kamu Tidak Bisa Events User Lain",
		})
		return
	}

	var input models.Event
	err := context.ShouldBindJSON(&input)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	config.DB.Model(&event).Updates(input)
	context.JSON(http.StatusOK, gin.H{
		"message": "Data Berhasil Di Update",
		"event":   event,
	})
}

func DeleteEvent(context *gin.Context) {
	userID, _ := context.Get("userID")

	var event models.Event
	paramsId := context.Param("id")

	var eventData = config.DB.First(&event, paramsId).Error
	if eventData != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "Data Event Tidak Ditemukan",
		})
		return
	}

	if event.UserID != uint(userID.(int)) {
		context.JSON(http.StatusForbidden, gin.H{
			"error": "Kamu Tidak Bisa Menghapus Event User Lain",
		})
		return
	}

	config.DB.Unscoped().Delete(&event)
	context.JSON(http.StatusOK, gin.H{
		"message": "Data Berhasil Di Delete",
		"event":   event,
	})
}
