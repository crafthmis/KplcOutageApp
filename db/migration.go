package db

import (
	"errors"
	"kplc-outage-app/models"
)

func Migrate() error {
	// Check if DB is initialized
	DB := GetDB()
	if DB == nil {
		return errors.New("database not initialized")
	}

	// Perform migrations
	err := DB.AutoMigrate(
		&models.Area{},
		&models.Contact{},
		&models.Constituency{},
		&models.County{},
		&models.Region{},
		&models.Campaign{},
		&models.Outage{},
		&models.OutageArea{},
		// Add other models as needed
	)

	if err != nil {
		return err
	}

	return nil
}
