package controllers

import (
	"fmt"
	"kplc-outage-app/models"
	"kplc-outage-app/services"
	"net/http"
	"sync"
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

	// Create campaigns concurrently
	errChan := make(chan error, len(input.Areas))
	var wg sync.WaitGroup

	for _, areaID := range input.Areas {
		wg.Add(1)
		go func(areaID uint) {
			defer wg.Done()
			err := processArea(areaID, outage.OtsID, input.Message)
			if err != nil {
				errChan <- err
			}
		}(areaID)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	var area models.Area
	err := services.GetAreaContactsByID(&area, fmt.Sprintf("%d", areaID))
	if err != nil {
		return fmt.Errorf("failed to fetch contacts for area %d: %v", areaID, err)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(area.Contacts))

	for _, contact := range area.Contacts {
		wg.Add(1)
		go func(contact models.Contact) {
			defer wg.Done()
			campaign := models.Campaign{
				OtsID:   otsID,
				Msisdn:  contact.Msisdn,
				Message: message,
			}
			err := services.CreateCampaign(&campaign)
			if err != nil {
				errChan <- fmt.Errorf("failed to save campaign for contact %s: %v", contact.Msisdn, err)
			}
		}(contact)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// func CreateOutageWithAreasHandler(c *gin.Context) {
// 	var input OutageInput

// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	outage := models.Outage{
// 		Message:    input.Message,
// 		OutageDate: input.OutageDate,
// 		SentStatus: input.SentStatus,
// 	}

// 	err := services.CreateOutageWithAreas(&outage, input.Areas)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create outage with areas"})
// 		return
// 	}

// 	//create campaign messages
// 	for _, areaId := range input.Areas {
// 		var area models.Area

// 		err = services.GetAreaContactsByID(&area, fmt.Sprintf("%d", areaId))
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contacts"})
// 			return
// 		}

// 		for _, contact := range area.Contacts {

// 			campaign := models.Campaign{
// 				OtsID:   outage.OtsID,
// 				Msisdn:  contact.Msisdn,
// 				Message: input.Message,
// 			}
// 			err := services.CreateCampaign(&campaign)
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save campaigns"})
// 				return
// 			}
// 		}

// 	}

// 	// Fetch the created outage with its areas
// 	createdOutage, err := services.GetOutageWithAreas(outage.OtsID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch created outage"})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{"outage": createdOutage})
// }
