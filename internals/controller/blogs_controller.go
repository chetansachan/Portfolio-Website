package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"portfolio-website/internals/database"
	"portfolio-website/internals/model"
	"portfolio-website/internals/service"
	"time"

	"github.com/google/uuid"
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

	ctx := context.Background()
	err := service.PostBlog_Service(database.Pool, ctx, b.Blog_Title, b.Blog_DES, b.Blog_Slug, b.Blog_Content, b.Blog_Banner)
	if err != nil {
		log.Println(err)
		return
	}

	blog := model.Blog{
		Blog_ID:      uuid.New(),
		Blog_Title:   b.Blog_Title,
		Blog_DES:     b.Blog_DES,
		Blog_Content: b.Blog_Content,
		Blog_Slug:    b.Blog_Slug,
		Blog_Banner:  b.Blog_Banner,
		Created_At:   time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(blog); err != nil {
		log.Println("Error Encoding JSON")
		return
	}
}
