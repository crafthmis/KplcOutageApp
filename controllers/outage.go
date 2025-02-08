package controllers

import (
	"fmt"
	"kplc-outage-app/models"
	"kplc-outage-app/services"
	"kplc-outage-app/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AreaInput struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
}

type OutageInput struct {
	Message    string      `json:"message"`
	OutageType string      `json:"outage_type"`
	OutageDate time.Time   `json:"outage_date"`
	SentStatus string      `json:"sent_status"`
	Areas      []AreaInput `json:"areas"`
}

func CreateOutageWithAreasHandler(c *gin.Context) {
	var input OutageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashInput := fmt.Sprintf("%s%s%s", input.Message, input.OutageDate, input.OutageType)
	myHash := utils.GenerateCryptoHash(hashInput)

	exists, err := services.GetOutageByHash(myHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Outage already exists"})
		return
	}

	outage := models.Outage{
		Message:    input.Message,
		OutageDate: input.OutageDate,
		OutageType: input.OutageType,
		SentStatus: input.SentStatus,
		OutageHash: myHash,
	}

	areaIDs := make([]uint, len(input.Areas))
	areaMsgs := make([]string, len(input.Areas))

	for i, area := range input.Areas {
		areaIDs[i] = area.Code
		areaMsgs[i] = area.Message
	}

	err2 := services.CreateOutageWithAreas(&outage, areaIDs, areaMsgs)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create outage with areas"})
		return
	}

	// Create campaigns concurrently
	for _, area := range input.Areas {
		err := processArea(area.Code, outage.OtsID, area.Message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create campaign"})
			return
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

func processArea(areaID uint, otsID uint, message string) error {
	var contacts []models.Contact
	err := services.GetAreaContactsByID(&contacts, fmt.Sprintf("%d", areaID))
	if err != nil {
		return fmt.Errorf("failed to fetch contacts for area %d: %v", areaID, err)
	}

	for _, contact := range contacts {
		fmt.Println("Contact " + contact.Msisdn)

		campaign := models.Campaign{
			OtsID:   otsID,
			Msisdn:  contact.Msisdn,
			Message: message,
		}
		err := services.CreateCampaign(&campaign)
		if err != nil {
			return fmt.Errorf("failed to create campaign for area %d: %v", areaID, err)
		}
	}
	return nil
}
