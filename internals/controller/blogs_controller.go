package controller

import (
	"encoding/json"
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
		return
	}

	var b model.Blog
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		log.Println("Error Decoding JSON")
		return
	}

	blog, err := c.BlogService.PostBlog_Service(r.Context(), b.Blog_Title, b.Blog_DES, b.Blog_Slug, b.Blog_Content, b.Blog_Banner)
	if err != nil {
		log.Println(err)
		return
	}

	response := model.PostBlogResponse{
		Response: model.Response{
			Success: true,
			Message: "Blog Created Successfully",
		},
		Data: model.Blog{
			Blog_ID:      blog.Blog_ID,
			Blog_Title:   blog.Blog_Title,
			Blog_DES:     blog.Blog_DES,
			Blog_Content: blog.Blog_Content,
			Blog_Slug:    blog.Blog_Slug,
			Blog_Banner:  blog.Blog_Banner,
			Created_At:   blog.Created_At,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error Encoding JSON")
		return
	}
}

func (c *BlogController) GetBlogBySlug_Controller(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Fatal("Invalid Method")
		return
	}

}
