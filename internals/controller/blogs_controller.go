package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"portfolio-website/internals/model"

	"github.com/google/uuid"
)

// BlogServiceInterface defines the contract for blog business logic.
// The controller depends on this interface (not the concrete service),
// enabling unit testing with mocks.
type BlogServiceInterface interface {
	PostBlog_Service(ctx context.Context, title, des, slug, content, banner string) (model.Blog, error)
	GetAllBlogs_Service(ctx context.Context) ([]model.Blog, error)
	GetBlogBySlug_Service(ctx context.Context, slug string) (model.Blog, error)
	GetBlogById_Service(ctx context.Context, id uuid.UUID) (model.Blog, error)
	DeleteBlogById_Service(ctx context.Context, id uuid.UUID) error
	UpdateBlogById_Service(ctx context.Context, id uuid.UUID, title, des, slug, content, banner string) (model.Blog, error)
}

type BlogController struct {
	BlogService BlogServiceInterface
}

func NewBlogController(blogService BlogServiceInterface) *BlogController {
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

		case errors.Is(err, model.ErrBlogAlreadyExist):
			w.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Blog Already exist",
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

func (c *BlogController) GetAllBlogs_Controller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.Response{Success: false, Message: "Method not allowed"})
		return
	}

	blogs, err := c.BlogService.GetAllBlogs_Service(r.Context())
	if err != nil {
		switch {
		case errors.Is(err, model.ErrBlogNotFound):
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "No Blog Found",
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
			Message: "Blogs Found Successfully",
		},
		Data: blogs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(response)
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

		case errors.Is(err, model.ErrBlogNotFound):
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Blog Not Found",
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
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(response)
}

func (c *BlogController) GetBlogById_Controller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.Response{Success: false, Message: "Method not allowed"})
		return
	}

	myid := r.PathValue("id")
	id, err := uuid.Parse(myid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.Response{
			Success: false,
			Message: "Invalid ID format",
		})
		return
	}

	blog, err := c.BlogService.GetBlogById_Service(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrIdValidation):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Valid ID is required",
			})
		case errors.Is(err, model.ErrBlogNotFound):
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Blog not found",
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(model.PostBlogResponse{
		Response: model.Response{
			Success: true,
			Message: "Blog Found Successfully",
		},
		Data: blog,
	})
}

func (c *BlogController) DeleteBlogById_Controller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.Response{Success: false, Message: "Method not allowed"})
		return
	}

	myid := r.PathValue("id")
	id, err := uuid.Parse(myid)
	fmt.Println(myid)
	fmt.Println(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.Response{
			Success: false,
			Message: "Invalid ID format",
		})
		return
	}

	if err := c.BlogService.DeleteBlogById_Service(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, model.ErrBlogNotFound):
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Blog Not Found",
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.Response{
		Success: true,
		Message: "Blog Deleted Successfully",
	})
}

func (c *BlogController) UpdateBlogById_Controller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.Response{Success: false, Message: "Method not allowed"})
		return
	}

	myid := r.PathValue("id")
	id, err := uuid.Parse(myid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.Response{
			Success: false,
			Message: "Invalid ID format",
		})
		return
	}

	var b model.Blog
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		log.Println("Error Decoding JSON")
		return
	}

	blog, err := c.BlogService.UpdateBlogById_Service(r.Context(), id, b.Blog_Title, b.Blog_DES, b.Blog_Slug, b.Blog_Content, b.Blog_Banner)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrIdValidation):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Id field is required",
			})

		case errors.Is(err, model.ErrBlogNotFound):
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Blog Not Found",
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.PostBlogResponse{
		Response: model.Response{
			Success: true,
			Message: "Blog Updated Successfully",
		},
		Data: blog,
	})
}
