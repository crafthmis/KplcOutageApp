package controllers

import (
	"fmt"
	"kplc-outage-app/models"
	"kplc-outage-app/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get Subscriptions
func GetSubscriptions(c *gin.Context) {
	var Subscription []models.Subscription

	err := services.GetAllSubscriptions(&Subscription)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, Subscription)
	}
}

// Get Subscription
func GetSubscription(c *gin.Context) {
	var Subscription models.Subscription
	id := c.Params.ByName("id")

	err := services.GetSubscriptionByID(&Subscription, id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, Subscription)
	}
}

// Create Subscription
func CreateSubscription(c *gin.Context) {
	var Subscription models.Subscription
	c.BindJSON(&Subscription)

	err := services.CreateSubscription(&Subscription)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, Subscription)
	}
}

// Update Subscription
func UpdateSubscription(c *gin.Context) {
	var Subscription models.Subscription
	id := c.Params.ByName("id")

	err := services.GetSubscriptionByID(&Subscription, id)
	if err != nil {
		c.JSON(http.StatusNotFound, Subscription)
	}

	c.BindJSON(&Subscription)

	err = services.UpdateSubscription(&Subscription, id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, Subscription)
	}
}

// Delete Subscription
func DeleteSubscription(c *gin.Context) {
	var Subscription models.Subscription
	id := c.Params.ByName("id")

	err := services.DeleteSubscription(&Subscription, id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, gin.H{"id" + id: "is deleted"})
	}
}
