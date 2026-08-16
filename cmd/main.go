package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"portfolio-website/internals/config"
	"portfolio-website/internals/database"
	"portfolio-website/internals/router"
	"portfolio-website/internals/scripts"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()
	fmt.Println(cfg)
	
	if err := database.Connect(ctx, cfg); err != nil {
		fmt.Println("Error Connecting to database")
		return
	}
	fmt.Println("Successfully connected to database")
	router := router.Router()

	scripts.RunMigration("file://./migration", cfg)

	fmt.Println("Server is running on http://localhost:8080...")
	err2 := http.ListenAndServe(":8080", router)
	if err2 != nil {
		log.Fatal("Server failed to start: ", err2)
	}
}
