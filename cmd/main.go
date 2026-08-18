package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is running on port %s...\n", port)
	err2 := http.ListenAndServe(":"+port, router)
	if err2 != nil {
		log.Fatal("Server failed to start: ", err2)
	}
}
