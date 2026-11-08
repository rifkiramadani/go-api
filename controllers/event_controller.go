package controllers

import (
	"context"
	"net/http"
	"os"
	"time"

	"example.com/event-app/config"
	"example.com/event-app/models"
	"github.com/gin-gonic/gin"
	"github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
)

func initImageKit() *imagekit.Client {
	client := imagekit.NewClient(
		option.WithPrivateKey(os.Getenv("IMAGEKIT_PRIVATE_KEY")),
	)
	return &client
}

func CreateEvent(c *gin.Context) {
	userID, _ := c.Get("userID")

	// Menerima File Form Data
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Gambar Wajib Di Upload",
		})
		return
	}
	defer file.Close()

	// 1. Upload File Ke Image Kit
	fileName := header.Filename
	ik := initImageKit()
	uploadRes, errUpload := ik.Files.Upload(context.Background(), imagekit.FileUploadParams{
		File:     file,
		FileName: fileName,
	})

	if errUpload != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal Upload Gambar Image Kit",
		})
		return
	}

	// created at
	parsedTime, _ := time.Parse(time.RFC3339, c.PostForm("datetime"))

	// 2. Simpan Ke Database
	event := models.Event{
		Name:        c.PostForm("name"),
		Description: c.PostForm("Description"),
		Location:    c.PostForm("location"),
		DateTime:    parsedTime,
		Image:       uploadRes.URL,
		ImageID:     uploadRes.FileID,
		UserID:      uint(userID.(int)),
	}

	config.DB.Create(&event)
	c.JSON(http.StatusCreated, gin.H{
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

func UpdateEvent(c *gin.Context) {
	userID, _ := c.Get("userID")

	var event models.Event
	paramsId := c.Param("id")

	var eventData = config.DB.First(&event, paramsId).Error
	if eventData != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Event Tidak Ditemukan",
		})
		return
	}

	if event.UserID != uint(userID.(int)) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Kamu Tidak Bisa Events User Lain",
		})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()

		ik := initImageKit()

		// upload file gambar baru
		fileName := header.Filename

		uploadRes, errUpload := ik.Files.Upload(context.Background(), imagekit.FileUploadParams{
			File:     file,
			FileName: fileName,
		})

		if errUpload == nil {
			// Hapus Gambar Lama
			if event.ImageID != "" {
				ik.Files.Delete(context.Background(), event.ImageID)
			}

			// Upload field gambar
			event.Image = uploadRes.URL
			event.ImageID = uploadRes.FileID
		}
	}

	// update data
	if name := c.PostForm("name"); name != "" {
		event.Name = name
	}
	if description := c.PostForm("description"); description != "" {
		event.Description = description
	}
	if location := c.PostForm("location"); location != "" {
		event.Location = location
	}
	if dateTimeStr := c.PostForm("datetime"); dateTimeStr != "" {
		parseTime, errParse := time.Parse(time.RFC3339, dateTimeStr)
		if errParse == nil {
			event.DateTime = parseTime
		}
	}

	config.DB.Save(&event)
	c.JSON(http.StatusOK, gin.H{
		"message": "Data Berhasil Di Update",
		"event":   event,
	})
}

func DeleteEvent(c *gin.Context) {
	userID, _ := c.Get("userID")

	var event models.Event
	paramsId := c.Param("id")

	var eventData = config.DB.First(&event, paramsId).Error
	if eventData != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Data Event Tidak Ditemukan",
		})
		return
	}

	if event.UserID != uint(userID.(int)) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Kamu Tidak Bisa Menghapus Event User Lain",
		})
		return
	}

	// hapus file image
	if event.ImageID != "" {
		ik := initImageKit()
		ik.Files.Delete(context.Background(), event.ImageID)
	}

	config.DB.Unscoped().Delete(&event)
	c.JSON(http.StatusOK, gin.H{
		"message": "Data Berhasil Di Delete",
		"event":   event,
	})
}
