package services

import (
	"kplc-outage-app/db"
	"kplc-outage-app/models"

	"gorm.io/gorm"
)

// CreateOutageWithAreas creates an Outage and its associated OutageAreas
func CreateOutageWithAreas(outage *models.Outage, areaIDs []uint, areaMessages []string) error {
	return db.GetDB().Transaction(func(tx *gorm.DB) error {
		// Create the outage
		if err := tx.Create(outage).Error; err != nil {
			db.GetDB().Rollback()
			return err
		}

		// Create outage areas
		for i, areaID := range areaIDs {
			outageArea := models.OutageArea{
				AreaID:  areaID,
				OtsID:   outage.OtsID,
				Message: areaMessages[i],
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

func GetOutageByHash(id string) (bool, error) {
	var count int64
	var outage models.Outage

	err := db.GetDB().Model(&outage).Where("outage_hash = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
