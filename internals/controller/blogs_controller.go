package controller

import (
	"log"
	"net/http"
)

func PostBlog_Controller(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Fatal("Invalid Method")
	}

}
