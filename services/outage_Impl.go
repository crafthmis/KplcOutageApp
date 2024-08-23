package services

import (
	"kplc-outage-app/db"
	"kplc-outage-app/models"

	"gorm.io/gorm"
)

// CreateOutageWithAreas creates an Outage and its associated OutageAreas
func CreateOutageWithAreas(outage *models.Outage, areaIDs []uint) error {
	return db.GetDB().Transaction(func(tx *gorm.DB) error {
		// Create the outage
		if err := tx.Create(outage).Error; err != nil {
			db.GetDB().Rollback()
			return err
		}

		// Create outage areas
		for _, areaID := range areaIDs {
			outageArea := models.OutageArea{
				AreaID: areaID,
				OtsID:  outage.OtsID,
			}
			if err := tx.Create(&outageArea).Error; err != nil {
				db.GetDB().Rollback()
				return err
			}
		}

		return nil
	})
}

func GetOutageWithAreas(outageID uint) (models.Outage, error) {
	var outage models.Outage
	err := db.GetDB().Preload("Areas").First(&outage, outageID).Error
	return outage, err
}
