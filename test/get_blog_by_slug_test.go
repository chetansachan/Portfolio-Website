package test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"portfolio-website/internals/controller"
	"portfolio-website/internals/model"

	"github.com/google/uuid"
)

// =============================================================================
// Tests for GET /blog/getblog/slug/{slug}
// =============================================================================

func TestGetBlogBySlug_Controller(t *testing.T) {
	sampleBlog := model.Blog{
		Blog_ID:      uuid.New(),
		Blog_Title:   "Mastering Go Concurrency",
		Blog_DES:     "Learn goroutines and channels in depth",
		Blog_Slug:    "mastering-go-concurrency",
		Blog_Content: "Detailed article on channels and select...",
		Blog_Banner:  "https://example.com/banner.png",
		Created_At:   time.Now(),
	}

	tests := []struct {
		name string

		// Request configuration
		method   string
		pathSlug string

		// Mock configuration
		getBySlugFunc func(ctx context.Context, slug string) (model.Blog, error)

		// Expected response
		expectedStatus  int
		expectedSuccess bool
		expectedMessage string
		expectData      bool
	}{
		// =====================================================================
		// HAPPY PATH
		// =====================================================================
		{
			name:     "Success — blog found by slug",
			method:   http.MethodGet,
			pathSlug: "mastering-go-concurrency",
			getBySlugFunc: func(ctx context.Context, slug string) (model.Blog, error) {
				return sampleBlog, nil
			},
			expectedStatus:  http.StatusFound,
			expectedSuccess: true,
			expectedMessage: "Blog Found Successfully",
			expectData:      true,
		},

		// =====================================================================
		// WRONG HTTP METHOD CASES
		// =====================================================================
		{
			name:            "Method Not Allowed — POST instead of GET",
			method:          http.MethodPost,
			pathSlug:        "mastering-go-concurrency",
			getBySlugFunc:   nil,
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedSuccess: false,
			expectedMessage: "Method not allowed",
		},
		{
			name:            "Method Not Allowed — DELETE instead of GET",
			method:          http.MethodDelete,
			pathSlug:        "mastering-go-concurrency",
			getBySlugFunc:   nil,
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedSuccess: false,
			expectedMessage: "Method not allowed",
		},
		{
			name:            "Method Not Allowed — PATCH instead of GET",
			method:          http.MethodPatch,
			pathSlug:        "mastering-go-concurrency",
			getBySlugFunc:   nil,
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedSuccess: false,
			expectedMessage: "Method not allowed",
		},

		// =====================================================================
		// VALIDATION / NOT FOUND CASES
		// =====================================================================
		{
			name:     "Validation Error — empty slug",
			method:   http.MethodGet,
			pathSlug: "",
			getBySlugFunc: func(ctx context.Context, slug string) (model.Blog, error) {
				return model.Blog{}, model.ErrSlugValidation
			},
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Slug field is required",
		},
		{
			name:     "Not Found — blog slug does not exist",
			method:   http.MethodGet,
			pathSlug: "non-existent-blog-slug",
			getBySlugFunc: func(ctx context.Context, slug string) (model.Blog, error) {
				return model.Blog{}, model.ErrBlogNotFound
			},
			expectedStatus:  http.StatusNotFound,
			expectedSuccess: false,
			expectedMessage: "Blog Not Found",
		},

		// =====================================================================
		// INTERNAL SERVER ERROR
		// =====================================================================
		{
			name:     "Internal Server Error — database failure",
			method:   http.MethodGet,
			pathSlug: "mastering-go-concurrency",
			getBySlugFunc: func(ctx context.Context, slug string) (model.Blog, error) {
				return model.Blog{}, errors.New("database connection lost")
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedSuccess: false,
			expectedMessage: "Internal server error. Please try again later.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockBlogService{
				GetBlogBySlugFunc: tt.getBySlugFunc,
			}
			blogController := controller.NewBlogController(mockService)

			req := httptest.NewRequest(tt.method, "/blog/getblog/slug/test", nil)
			req.SetPathValue("slug", tt.pathSlug)
			w := httptest.NewRecorder()

			blogController.GetBlogBySlug_Controller(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status code: got %d, want %d", w.Code, tt.expectedStatus)
			}

			if w.Header().Get("Content-Type") != "application/json" {
				t.Errorf("Content-Type: got %q, want %q", w.Header().Get("Content-Type"), "application/json")
			}

			if tt.expectData {
				var resp model.GetBlogsResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if resp.Success != tt.expectedSuccess {
					t.Errorf("Success: got %v, want %v", resp.Success, tt.expectedSuccess)
				}
				if resp.Message != tt.expectedMessage {
					t.Errorf("Message: got %q, want %q", resp.Message, tt.expectedMessage)
				}
				if len(resp.Data) != 1 || resp.Data[0].Blog_Slug != tt.pathSlug {
					t.Errorf("Data: unexpected payload %+v", resp.Data)
				}
			} else {
				var resp model.Response
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("Failed to decode error response: %v", err)
				}
				if resp.Success != tt.expectedSuccess {
					t.Errorf("Success: got %v, want %v", resp.Success, tt.expectedSuccess)
				}
				if resp.Message != tt.expectedMessage {
					t.Errorf("Message: got %q, want %q", resp.Message, tt.expectedMessage)
				}
			}
		})
	}
}
