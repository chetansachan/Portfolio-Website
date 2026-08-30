package test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"portfolio-website/internals/controller"
	"portfolio-website/internals/model"

	"github.com/google/uuid"
)

// =============================================================================
// Tests for POST /blog/postblog
// =============================================================================

func TestPostBlog_Controller(t *testing.T) {
	sampleBlog := model.Blog{
		Blog_ID:      uuid.New(),
		Blog_Title:   "Unit Testing in Go",
		Blog_DES:     "A comprehensive guide to testing in Go",
		Blog_Slug:    "unit-testing-in-go",
		Blog_Content: "Full article content here...",
		Blog_Banner:  "https://example.com/banner.png",
		Created_At:   time.Now(),
	}

	tests := []struct {
		name string // descriptive test case name

		// Request configuration
		method      string // HTTP method to use
		contentType string // Content-Type header (empty = default)
		body        string // raw JSON body

		// Mock configuration
		postFunc func(ctx context.Context, title, des, slug, content, banner string) (model.Blog, error)

		// Expected response
		expectedStatus  int
		expectedSuccess bool
		expectedMessage string
		expectData      bool // whether response should include blog data
	}{
		// =====================================================================
		// HAPPY PATH
		// =====================================================================
		{
			name:   "Success — all fields provided, blog created",
			method: http.MethodPost,
			body: `{
				"blog_title": "Unit Testing in Go",
				"blog_des": "A comprehensive guide to testing in Go",
				"blog_slug": "unit-testing-in-go",
				"blog_content": "Full article content here...",
				"blog_banner": "https://example.com/banner.png"
			}`,
			postFunc: func(ctx context.Context, title, des, slug, content, banner string) (model.Blog, error) {
				return sampleBlog, nil
			},
			expectedStatus:  http.StatusCreated,
			expectedSuccess: true,
			expectedMessage: "Blog created successfully",
			expectData:      true,
		},

		// =====================================================================
		// WRONG HTTP METHOD CASES
		// =====================================================================
		{
			name:            "Method Not Allowed — GET instead of POST",
			method:          http.MethodGet,
			body:            "",
			postFunc:        nil,
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedSuccess: false,
			expectedMessage: "Method not allowed",
		},
		{
			name:            "Method Not Allowed — DELETE instead of POST",
			method:          http.MethodDelete,
			body:            "",
			postFunc:        nil,
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedSuccess: false,
			expectedMessage: "Method not allowed",
		},
		{
			name:            "Method Not Allowed — PUT instead of POST",
			method:          http.MethodPut,
			body:            "",
			postFunc:        nil,
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedSuccess: false,
			expectedMessage: "Method not allowed",
		},
		{
			name:            "Method Not Allowed — PATCH instead of POST",
			method:          http.MethodPatch,
			body:            "",
			postFunc:        nil,
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedSuccess: false,
			expectedMessage: "Method not allowed",
		},

		// =====================================================================
		// MALFORMED / INVALID JSON BODY
		// =====================================================================
		{
			name:            "Bad Request — completely invalid JSON (not JSON at all)",
			method:          http.MethodPost,
			body:            `this is not json`,
			postFunc:        nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Invalid request body",
		},
		{
			name:            "Bad Request — empty body (no JSON sent)",
			method:          http.MethodPost,
			body:            "",
			postFunc:        nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Invalid request body",
		},
		{
			name:            "Bad Request — partial/truncated JSON",
			method:          http.MethodPost,
			body:            `{"blog_title": "Hello`,
			postFunc:        nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Invalid request body",
		},
		{
			name:            "Bad Request — JSON array instead of object",
			method:          http.MethodPost,
			body:            `[{"blog_title": "Hello"}]`,
			postFunc:        nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Invalid request body",
		},

		// =====================================================================
		// VALIDATION FAILURES (empty required fields)
		// =====================================================================
		{
			name:   "Validation Error — all fields empty strings",
			method: http.MethodPost,
			body: `{
				"blog_title": "",
				"blog_des": "",
				"blog_slug": "",
				"blog_content": "",
				"blog_banner": ""
			}`,
			postFunc: func(ctx context.Context, title, des, slug, content, banner string) (model.Blog, error) {
				return model.Blog{}, model.ErrValidation
			},
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "All fields are required",
		},
		{
			name:   "Validation Error — missing title only",
			method: http.MethodPost,
			body: `{
				"blog_des": "Some desc",
				"blog_slug": "some-slug",
				"blog_content": "Some content",
				"blog_banner": "https://example.com/img.png"
			}`,
			postFunc: func(ctx context.Context, title, des, slug, content, banner string) (model.Blog, error) {
				return model.Blog{}, model.ErrValidation
			},
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "All fields are required",
		},
		{
			name:   "Validation Error — missing banner only",
			method: http.MethodPost,
			body: `{
				"blog_title": "My Blog",
				"blog_des": "Some desc",
				"blog_slug": "my-blog",
				"blog_content": "Full content"
			}`,
			postFunc: func(ctx context.Context, title, des, slug, content, banner string) (model.Blog, error) {
				return model.Blog{}, model.ErrValidation
			},
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "All fields are required",
		},
		{
			name:   "Validation Error — empty JSON object (no fields at all)",
			method: http.MethodPost,
			body:   `{}`,
			postFunc: func(ctx context.Context, title, des, slug, content, banner string) (model.Blog, error) {
				return model.Blog{}, model.ErrValidation
			},
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "All fields are required",
		},

		// =====================================================================
		// DUPLICATE / CONFLICT
		// =====================================================================
		{
			name:   "Conflict — blog with same slug already exists",
			method: http.MethodPost,
			body: `{
				"blog_title": "Unit Testing in Go",
				"blog_des": "A comprehensive guide",
				"blog_slug": "unit-testing-in-go",
				"blog_content": "Full content...",
				"blog_banner": "https://example.com/banner.png"
			}`,
			postFunc: func(ctx context.Context, title, des, slug, content, banner string) (model.Blog, error) {
				return sampleBlog, model.ErrBlogAlreadyExist
			},
			expectedStatus:  http.StatusConflict,
			expectedSuccess: false,
			expectedMessage: "Blog Already exist",
		},

		// =====================================================================
		// INTERNAL SERVER ERROR (unexpected failures)
		// =====================================================================
		{
			name:   "Internal Server Error — database write failed",
			method: http.MethodPost,
			body: `{
				"blog_title": "My Blog",
				"blog_des": "Description",
				"blog_slug": "my-blog",
				"blog_content": "Content",
				"blog_banner": "https://example.com/b.png"
			}`,
			postFunc: func(ctx context.Context, title, des, slug, content, banner string) (model.Blog, error) {
				return model.Blog{}, errors.New("connection refused")
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedSuccess: false,
			expectedMessage: "Internal server error. Please try again later.",
		},
		{
			name:   "Internal Server Error — unexpected unknown error from service",
			method: http.MethodPost,
			body: `{
				"blog_title": "My Blog",
				"blog_des": "Description",
				"blog_slug": "my-blog",
				"blog_content": "Content",
				"blog_banner": "https://example.com/b.png"
			}`,
			postFunc: func(ctx context.Context, title, des, slug, content, banner string) (model.Blog, error) {
				return model.Blog{}, errors.New("disk full: cannot write WAL")
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedSuccess: false,
			expectedMessage: "Internal server error. Please try again later.",
		},

		// =====================================================================
		// EDGE CASES — unusual but valid/invalid inputs
		// =====================================================================
		{
			name:   "Edge Case — extra unknown fields in JSON body (should be ignored)",
			method: http.MethodPost,
			body: `{
				"blog_title": "My Blog",
				"blog_des": "Description",
				"blog_slug": "my-blog",
				"blog_content": "Content",
				"blog_banner": "https://example.com/b.png",
				"unknown_field": "should be ignored",
				"another_field": 12345
			}`,
			postFunc: func(ctx context.Context, title, des, slug, content, banner string) (model.Blog, error) {
				return sampleBlog, nil
			},
			expectedStatus:  http.StatusCreated,
			expectedSuccess: true,
			expectedMessage: "Blog created successfully",
			expectData:      true,
		},
		{
			name:   "Edge Case — very long title (stress test)",
			method: http.MethodPost,
			body: `{
				"blog_title": "` + strings.Repeat("A", 5000) + `",
				"blog_des": "Description",
				"blog_slug": "very-long-title-blog",
				"blog_content": "Content",
				"blog_banner": "https://example.com/b.png"
			}`,
			postFunc: func(ctx context.Context, title, des, slug, content, banner string) (model.Blog, error) {
				return sampleBlog, nil
			},
			expectedStatus:  http.StatusCreated,
			expectedSuccess: true,
			expectedMessage: "Blog created successfully",
			expectData:      true,
		},
		{
			name:   "Edge Case — special characters in slug",
			method: http.MethodPost,
			body: `{
				"blog_title": "Test <script>alert('xss')</script>",
				"blog_des": "Description with \"quotes\" and emoji 🚀",
				"blog_slug": "test-special-chars-<>",
				"blog_content": "Content with\nnewlines\tand\ttabs",
				"blog_banner": "https://example.com/b.png"
			}`,
			postFunc: func(ctx context.Context, title, des, slug, content, banner string) (model.Blog, error) {
				return sampleBlog, nil
			},
			expectedStatus:  http.StatusCreated,
			expectedSuccess: true,
			expectedMessage: "Blog created successfully",
			expectData:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// --- Setup mock ---
			mockService := &MockBlogService{
				PostFunc: func(ctx context.Context, title, des, slug, content, banner string, featured bool) (model.Blog, error) {
					if tt.postFunc == nil {
						return model.Blog{}, nil
					}
					return tt.postFunc(ctx, title, des, slug, content, banner)
				},
			}
			blogController := controller.NewBlogController(mockService)

			// --- Build request ---
			var body *strings.Reader
			if tt.body != "" {
				body = strings.NewReader(tt.body)
			} else {
				body = strings.NewReader("")
			}
			req := httptest.NewRequest(tt.method, "/blog/postblog", body)
			req.Header.Set("Content-Type", "application/json")

			// --- Record response ---
			w := httptest.NewRecorder()
			blogController.CreateBlog_Controller(w, req)

			// --- Assert status code ---
			if w.Code != tt.expectedStatus {
				t.Errorf("Status code: got %d, want %d", w.Code, tt.expectedStatus)
			}

			// --- Assert Content-Type header ---
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Content-Type: got %q, want %q", contentType, "application/json")
			}

			// --- Assert response body ---
			if tt.expectData {
				// For success cases, decode as PostBlogResponse (includes Data field)
				var resp model.PostBlogResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}

				if resp.Success != tt.expectedSuccess {
					t.Errorf("Success: got %v, want %v", resp.Success, tt.expectedSuccess)
				}
				if resp.Message != tt.expectedMessage {
					t.Errorf("Message: got %q, want %q", resp.Message, tt.expectedMessage)
				}
				// Verify blog data is populated
				if resp.Data.Blog_ID == uuid.Nil {
					t.Error("Expected blog data with a valid Blog_ID, got nil UUID")
				}
			} else {
				// For error cases, decode as simple Response
				var resp model.Response
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
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
