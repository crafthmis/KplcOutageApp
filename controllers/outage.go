package controllers

import (
	"kplc-outage-app/models"
	"kplc-outage-app/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type OutageInput struct {
	Message    string    `json:"message"`
	OutageDate time.Time `json:"outage_date"`
	SendStatus string    `json:"send_status"`
	Areas      []uint    `json:"areas"` // List of area IDs
}

func CreateOutageWithAreasHandler(c *gin.Context) {
	var input OutageInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	outage := models.Outage{
		Message:    input.Message,
		OutageDate: input.OutageDate,
		SendStatus: input.SendStatus,
	}

	err := services.CreateOutageWithAreas(&outage, input.Areas)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create outage with areas"})
		return
	}

	// Fetch the created outage with its areas
	createdOutage, err := services.GetOutageWithAreas(outage.OtsID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch created outage"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"outage": createdOutage})
}
