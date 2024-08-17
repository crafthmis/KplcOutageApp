package main

import (
	"fmt"
	"kplc-outage-app/config"
	"kplc-outage-app/db"
	"kplc-outage-app/routes"
	"os"
)

func main() {
	err := config.EnvSetup()
	if err != nil {
		fmt.Printf("Error setting up environment: %v\n", err)
		return
	}
	fmt.Println("Environment variables loaded successfully")
	fmt.Printf("Current environment: %s\n", os.Getenv("GO_ENV"))
	err2 := db.DatabaseInit()
	if err2 != nil {
		fmt.Printf("Error in connecting to db: %v\n", err2)
		return
	}
	fmt.Println("DB successfully connected", os.Getenv("POSTGRES_DB"))
	
	r := routes.SetupRouter()
	//running
	r.Run()

}
