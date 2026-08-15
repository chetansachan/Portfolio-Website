package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"portfolio-website/internals/model"
)

func PostBlog_Controller(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Fatal("Invalid Method")
	}

	var b model.Blog
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		log.Println("Error Decoding JSON")
		return
	}

}
