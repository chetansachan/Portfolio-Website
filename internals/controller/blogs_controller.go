package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"portfolio-website/internals/model"
	"portfolio-website/internals/service"
)

type BlogController struct {
	BlogService *service.BlogService
}

func NewBlogController(blogService *service.BlogService) *BlogController {
	return &BlogController{BlogService: blogService}
}

func (c *BlogController) PostBlog_Controller(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Fatal("Invalid Method")
	}

	var b model.Blog
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		log.Println("Error Decoding JSON")
		return
	}

	ctx := context.Background()
	err := c.BlogService.PostBlog_Service(ctx, b.Blog_Title, b.Blog_DES, b.Blog_Slug, b.Blog_Content, b.Blog_Banner)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(b.Blog_ID)

	response := model.PostBlogResponse{
		Response: model.Response{
			Success: true,
			Message: "Blog Created Successfully",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error Encoding JSON")
		return
	}
}
