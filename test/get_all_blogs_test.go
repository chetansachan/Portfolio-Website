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
// Tests for GET /blog/getallblogs
// =============================================================================

func TestGetAllBlogs_Controller(t *testing.T) {
	sampleBlogs := []model.Blog{
		{
			Blog_ID:      uuid.New(),
			Blog_Title:   "First Blog",
			Blog_DES:     "First Description",
			Blog_Slug:    "first-blog",
			Blog_Content: "Content 1",
			Blog_Banner:  "https://example.com/b1.png",
			Created_At:   time.Now(),
		},
		{
			Blog_ID:      uuid.New(),
			Blog_Title:   "Second Blog",
			Blog_DES:     "Second Description",
			Blog_Slug:    "second-blog",
			Blog_Content: "Content 2",
			Blog_Banner:  "https://example.com/b2.png",
			Created_At:   time.Now().Add(-time.Hour),
		},
	}

	tests := []struct {
		name string

		// Request configuration
		method      string
		queryParams string

		// Mock configuration
		getAllFunc func(ctx context.Context, size int, offset int, featured *bool) ([]model.Blog, int, error)

		// Expected response
		expectedStatus  int
		expectedSuccess bool
		expectedMessage string
		expectDataCount int
		expectTotal     int
		expectPage      int
		expectSize      int
	}{
		// =====================================================================
		// HAPPY PATH
		// =====================================================================
		{
			name:        "Success — default pagination params (page=1, size=10)",
			method:      http.MethodGet,
			queryParams: "",
			getAllFunc: func(ctx context.Context, size int, offset int, featured *bool) ([]model.Blog, int, error) {
				if size != 10 || offset != 0 {
					t.Errorf("Expected size=10, offset=0, got size=%d, offset=%d", size, offset)
				}
				return sampleBlogs, 2, nil
			},
			expectedStatus:  http.StatusOK,
			expectedSuccess: true,
			expectedMessage: "Blogs Found Successfully",
			expectDataCount: 2,
			expectTotal:     2,
			expectPage:      1,
			expectSize:      10,
		},
		{
			name:        "Success — custom valid pagination (page=2, size=5)",
			method:      http.MethodGet,
			queryParams: "?page=2&size=5",
			getAllFunc: func(ctx context.Context, size int, offset int, featured *bool) ([]model.Blog, int, error) {
				if size != 5 || offset != 5 {
					t.Errorf("Expected size=5, offset=5, got size=%d, offset=%d", size, offset)
				}
				return sampleBlogs[:1], 15, nil
			},
			expectedStatus:  http.StatusOK,
			expectedSuccess: true,
			expectedMessage: "Blogs Found Successfully",
			expectDataCount: 1,
			expectTotal:     15,
			expectPage:      2,
			expectSize:      5,
		},
		{
			name:        "Success — empty list when no blogs exist",
			method:      http.MethodGet,
			queryParams: "?page=1&size=10",
			getAllFunc: func(ctx context.Context, size int, offset int, featured *bool) ([]model.Blog, int, error) {
				return []model.Blog{}, 0, nil
			},
			expectedStatus:  http.StatusOK,
			expectedSuccess: true,
			expectedMessage: "Blogs Found Successfully",
			expectDataCount: 0,
			expectTotal:     0,
			expectPage:      1,
			expectSize:      10,
		},
		{
			name:        "Success — size capped at maximum of 100",
			method:      http.MethodGet,
			queryParams: "?page=1&size=500",
			getAllFunc: func(ctx context.Context, size int, offset int, featured *bool) ([]model.Blog, int, error) {
				if size != 100 {
					t.Errorf("Expected size to be capped at 100, got %d", size)
				}
				return sampleBlogs, 2, nil
			},
			expectedStatus:  http.StatusOK,
			expectedSuccess: true,
			expectedMessage: "Blogs Found Successfully",
			expectDataCount: 2,
			expectTotal:     2,
			expectPage:      1,
			expectSize:      100,
		},

		// =====================================================================
		// WRONG HTTP METHOD CASES
		// =====================================================================
		{
			name:            "Method Not Allowed — POST instead of GET",
			method:          http.MethodPost,
			queryParams:     "",
			getAllFunc:      nil,
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedSuccess: false,
			expectedMessage: "Method not allowed",
		},
		{
			name:            "Method Not Allowed — DELETE instead of GET",
			method:          http.MethodDelete,
			queryParams:     "",
			getAllFunc:      nil,
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedSuccess: false,
			expectedMessage: "Method not allowed",
		},
		{
			name:            "Method Not Allowed — PUT instead of GET",
			method:          http.MethodPut,
			queryParams:     "",
			getAllFunc:      nil,
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedSuccess: false,
			expectedMessage: "Method not allowed",
		},

		// =====================================================================
		// INVALID QUERY PARAMETERS
		// =====================================================================
		{
			name:            "Bad Request — invalid page format (non-integer)",
			method:          http.MethodGet,
			queryParams:     "?page=abc",
			getAllFunc:      nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Invalid 'page' parameter: must be a positive integer",
		},
		{
			name:            "Bad Request — zero page number",
			method:          http.MethodGet,
			queryParams:     "?page=0",
			getAllFunc:      nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Invalid 'page' parameter: must be a positive integer",
		},
		{
			name:            "Bad Request — negative page number",
			method:          http.MethodGet,
			queryParams:     "?page=-1",
			getAllFunc:      nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Invalid 'page' parameter: must be a positive integer",
		},
		{
			name:            "Bad Request — invalid size format (non-integer)",
			method:          http.MethodGet,
			queryParams:     "?size=xyz",
			getAllFunc:      nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Invalid 'size' parameter: must be a positive integer",
		},
		{
			name:            "Bad Request — zero size",
			method:          http.MethodGet,
			queryParams:     "?size=0",
			getAllFunc:      nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Invalid 'size' parameter: must be a positive integer",
		},
		{
			name:            "Bad Request — negative size",
			method:          http.MethodGet,
			queryParams:     "?size=-5",
			getAllFunc:      nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Invalid 'size' parameter: must be a positive integer",
		},

		// =====================================================================
		// NOT FOUND / SERVICE ERRORS
		// =====================================================================
		{
			name:        "Not Found — service returns ErrBlogNotFound",
			method:      http.MethodGet,
			queryParams: "",
			getAllFunc: func(ctx context.Context, size int, offset int, featured *bool) ([]model.Blog, int, error) {
				return nil, 0, model.ErrBlogNotFound
			},
			expectedStatus:  http.StatusNotFound,
			expectedSuccess: false,
			expectedMessage: "No Blog Found",
		},
		{
			name:        "Internal Server Error — database connection failure",
			method:      http.MethodGet,
			queryParams: "",
			getAllFunc: func(ctx context.Context, size int, offset int, featured *bool) ([]model.Blog, int, error) {
				return nil, 0, errors.New("db connection timeout")
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedSuccess: false,
			expectedMessage: "Internal server error. Please try again later.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockBlogService{
				GetAllBlogsFunc: tt.getAllFunc,
			}
			blogController := controller.NewBlogController(mockService)

			req := httptest.NewRequest(tt.method, "/blog/getallblogs"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			blogController.GetAllBlogs_Controller(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status code: got %d, want %d", w.Code, tt.expectedStatus)
			}

			if w.Header().Get("Content-Type") != "application/json" {
				t.Errorf("Content-Type: got %q, want %q", w.Header().Get("Content-Type"), "application/json")
			}

			if tt.expectedStatus == http.StatusOK {
				var resp model.PaginatedResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if resp.Success != tt.expectedSuccess {
					t.Errorf("Success: got %v, want %v", resp.Success, tt.expectedSuccess)
				}
				if resp.Message != tt.expectedMessage {
					t.Errorf("Message: got %q, want %q", resp.Message, tt.expectedMessage)
				}
				if len(resp.Data) != tt.expectDataCount {
					t.Errorf("Data count: got %d, want %d", len(resp.Data), tt.expectDataCount)
				}
				if resp.Meta.TotalBlogs != tt.expectTotal {
					t.Errorf("Meta TotalBlogs: got %d, want %d", resp.Meta.TotalBlogs, tt.expectTotal)
				}
				if resp.Meta.CurrentPage != tt.expectPage {
					t.Errorf("Meta CurrentPage: got %d, want %d", resp.Meta.CurrentPage, tt.expectPage)
				}
				if resp.Meta.PageSize != tt.expectSize {
					t.Errorf("Meta PageSize: got %d, want %d", resp.Meta.PageSize, tt.expectSize)
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
