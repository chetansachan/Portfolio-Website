package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"portfolio-website/internals/model"
	"portfolio-website/internals/service"

	"github.com/google/uuid"
)

type BlogController struct {
	BlogService *service.BlogService
}

func NewBlogController(blogService *service.BlogService) *BlogController {
	return &BlogController{BlogService: blogService}
}

func (c *BlogController) PostBlog_Controller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.Response{Success: false, Message: "Method not allowed"})
		return
	}

	var b model.Blog
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		log.Printf("[WARN] Invalid JSON payload: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.Response{Success: false, Message: "Invalid request body"})
		return
	}

	blog, err := c.BlogService.PostBlog_Service(r.Context(), b.Blog_Title, b.Blog_DES, b.Blog_Slug, b.Blog_Content, b.Blog_Banner)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrValidation):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "All fields are required",
			})

		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Internal server error. Please try again later.",
			})
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(model.PostBlogResponse{
		Response: model.Response{Success: true, Message: "Blog created successfully"},
		Data:     blog,
	})
}

func (c *BlogController) GetBlogBySlug_Controller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.Response{Success: false, Message: "Method not allowed"})
		return
	}

	slug := r.PathValue("slug")
	blog, err := c.BlogService.GetBlogBySlug_Service(r.Context(), slug)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrSlugValidation):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Slug field is required",
			})

		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Internal server error. Please try again later.",
			})
		}
		return
	}

	response := model.GetBlogsResponse{
		Response: model.Response{
			Success: true,
			Message: "Blog Found Successfully",
		},
		Data: []model.Blog{blog},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *BlogController) DeleteBlogById_Controller(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	myid := r.PathValue("id")
	id, err := uuid.Parse(myid)
	if err != nil {
		log.Println("Error Parsing string")
		return
	}

	if err := c.BlogService.DeleteBlogById_Service(r.Context(), id); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.Response{
			Success: false,
			Message: "Blog not found",
		})
		return
	}

	response := model.Response{
		Success: true,
		Message: "Blog Deleted Successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *BlogController) UpdateBlogById_Controller(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	myid := r.PathValue("id")
	id, err := uuid.Parse(myid)
	if err != nil {
		log.Println("Error Parsing string")
		return
	}

	var b model.Blog
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		log.Println("Error Decoding JSON")
		return
	}

	blog, err := c.BlogService.UpdateBlogById_Service(r.Context(), id, b.Blog_Title, b.Blog_DES, b.Blog_Slug, b.Blog_Content, b.Blog_Banner)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.Response{
			Success: false,
			Message: "Blog not found",
		})
		return
	}

	response := model.PostBlogResponse{
		Response: model.Response{
			Success: true,
			Message: "Blog Updated Successfully",
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
