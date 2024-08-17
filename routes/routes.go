package routes

import (
	"kplc-outage-app/controllers"
	"kplc-outage-app/ussd"

	"github.com/gin-gonic/gin"
)

// Setup Router
func SetupRouter() *gin.Engine {
	r := gin.Default()
	// grouping
	grp := r.Group("/api/v1")
	{
		grp.GET("/regions", controllers.GetRegions)
		grp.GET("/region/:id", controllers.GetRegion)
		grp.POST("/region", controllers.CreateRegion)
		grp.PUT("/region/:id", controllers.UpdateRegion)
		grp.DELETE("/region/:id", controllers.DeleteRegion)
		//counties
		grp.GET("/counties", controllers.GetCounties)
		grp.GET("/county/:id", controllers.GetCounty)
		grp.POST("/county", controllers.CreateCounty)
		grp.PUT("/county/:id", controllers.UpdateCounty)
		grp.DELETE("/county/:id", controllers.DeleteCounty)
		//constituencies
		grp.GET("/constituencies", controllers.GetConstituencies)
		grp.GET("/constituency/:id", controllers.GetConstituency)
		grp.POST("/constituency", controllers.CreateConstituency)
		grp.PUT("/constituency/:id", controllers.UpdateConstituency)
		grp.DELETE("/constituency/:id", controllers.DeleteConstituency)
		//area
		grp.GET("/areas", controllers.GetAreas)
		grp.GET("/area/:id", controllers.GetArea)
		grp.GET("/area/:id/contacts", controllers.GetAreaContacts)
		grp.POST("/area", controllers.CreateArea)
		grp.PUT("/area/:id", controllers.UpdateArea)
		grp.DELETE("/area/:id", controllers.DeleteArea)
		//subscription
		grp.GET("/subscriptions", controllers.GetSubscriptions)
		grp.GET("/subscription/:id", controllers.GetSubscription)
		grp.POST("/subscription", controllers.CreateSubscription)
		grp.PUT("/subscription/:id", controllers.UpdateSubscription)
		grp.DELETE("/subscription/:id", controllers.DeleteSubscription)

		// ussd callback
		grp.POST("/ussd_callback", gin.WrapF(ussd.UssdCallback))

	}
	return r
}
