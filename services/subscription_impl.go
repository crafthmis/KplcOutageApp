package services

import (
	"fmt"
	"kplc-outage-app/db"
	"kplc-outage-app/models"
)

func GetAllSubscriptions(Subscription *[]models.Subscription) (err error) {
	if err = db.GetDB().Find(Subscription).Error; err != nil {
		return err
	}
	return nil
}

// Create Subscription
func CreateSubscription(Subscription *models.Subscription) (err error) {
	if err = db.GetDB().Create(Subscription).Error; err != nil {
		return err
	}
	return nil
}

// Get Subscription ByID
func GetSubscriptionByID(Subscription *models.Subscription, id string) (err error) {
	if err = db.GetDB().Where("sub_id = ?", id).First(Subscription).Error; err != nil {
		return err
	}
	return nil
}

func GetSubscriptionContactsByID(Subscription *models.Subscription, id string) (err error) {
	if err = db.GetDB().Preload("Contacts").Where("sub_id = ?", id).First(Subscription).Error; err != nil {
		return err
	}
	return nil
}

// Update Subscription
func UpdateSubscription(Subscription *models.Subscription, id string) (err error) {
	fmt.Println(Subscription)
	db.GetDB().Save(Subscription)
	return nil
}

// Delete Subscription
func DeleteSubscription(Subscription *models.Subscription, id string) (err error) {
	db.GetDB().Where("sub_id = ?", id).Delete(Subscription)
	return nil
}
