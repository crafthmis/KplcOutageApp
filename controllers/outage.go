package controllers

import (
	"fmt"
	"kplc-outage-app/models"
	"kplc-outage-app/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type OutageInput struct {
	Message    string    `json:"message"`
	OutageDate time.Time `json:"outage_date"`
	SentStatus string    `json:"sent_status"`
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
		SentStatus: input.SentStatus,
	}

	err := services.CreateOutageWithAreas(&outage, input.Areas)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create outage with areas"})
		return
	}

	//create campaign messages
	for _, areaId := range input.Areas {
		var area models.Area

		err = services.GetAreaContactsByID(&area, fmt.Sprintf("%d", areaId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contacts"})
			return
		}

		for _, contact := range area.Contacts {

			campaign := models.Campaign{
				OtsID:   outage.OtsID,
				Msisdn:  contact.Msisdn,
				Message: input.Message,
			}
			err := services.CreateCampaign(&campaign)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save campaigns"})
				return
			}
		}

	}

	// Fetch the created outage with its areas
	createdOutage, err := services.GetOutageWithAreas(outage.OtsID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch created outage"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"outage": createdOutage})
}
