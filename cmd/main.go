package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"portfolio-website/internals/config"
	"portfolio-website/internals/database"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	_, err := database.Connect(ctx, cfg)
	if err != nil {
		fmt.Println("Error Connecting to database")
		return
	}
	fmt.Println("Successfully connected to database")

	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello")
	})



	fmt.Println("Server is running on http://localhost:8080...")
	err2 := http.ListenAndServe(":8080", nil)
	if err2 != nil {
		log.Fatal("Server failed to start: ", err2)
	}
}
