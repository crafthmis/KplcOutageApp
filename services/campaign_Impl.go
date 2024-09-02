package services

import (
	"kplc-outage-app/db"
	"kplc-outage-app/models"
)

// Create Campaign
func CreateCampaign(Campaign *models.Campaign) (err error) {
	if err = db.GetDB().Create(Campaign).Error; err != nil {
		return err
	}
	return nil
}
