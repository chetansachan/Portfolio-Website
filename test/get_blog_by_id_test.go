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
// Tests for GET /blog/getblog/id/{id}
// =============================================================================

func TestGetBlogById_Controller(t *testing.T) {
	testUUID := uuid.New()
	sampleBlog := model.Blog{
		Blog_ID:      testUUID,
		Blog_Title:   "Building Scalable APIs",
		Blog_DES:     "How to build high-performance APIs",
		Blog_Slug:    "building-scalable-apis",
		Blog_Content: "Deep dive into REST and HTTP/2 in Go",
		Blog_Banner:  "https://example.com/banner.png",
		Created_At:   time.Now(),
	}

	tests := []struct {
		name string

		// Request configuration
		method string
		pathID string

		// Mock configuration
		getByIdFunc func(ctx context.Context, id uuid.UUID) (model.Blog, error)

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
			name:   "Success — blog found by valid UUID",
			method: http.MethodGet,
			pathID: testUUID.String(),
			getByIdFunc: func(ctx context.Context, id uuid.UUID) (model.Blog, error) {
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
			pathID:          testUUID.String(),
			getByIdFunc:     nil,
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedSuccess: false,
			expectedMessage: "Method not allowed",
		},
		{
			name:            "Method Not Allowed — DELETE instead of GET",
			method:          http.MethodDelete,
			pathID:          testUUID.String(),
			getByIdFunc:     nil,
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedSuccess: false,
			expectedMessage: "Method not allowed",
		},

		// =====================================================================
		// INVALID ID FORMAT CASES
		// =====================================================================
		{
			name:            "Bad Request — non-UUID string",
			method:          http.MethodGet,
			pathID:          "invalid-uuid-format",
			getByIdFunc:     nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Invalid ID format",
		},
		{
			name:            "Bad Request — numeric ID",
			method:          http.MethodGet,
			pathID:          "12345",
			getByIdFunc:     nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Invalid ID format",
		},
		{
			name:            "Bad Request — empty ID string",
			method:          http.MethodGet,
			pathID:          "",
			getByIdFunc:     nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Invalid ID format",
		},

		// =====================================================================
		// NOT FOUND / VALIDATION ERRORS
		// =====================================================================
		{
			name:   "Not Found — blog does not exist",
			method: http.MethodGet,
			pathID: uuid.New().String(),
			getByIdFunc: func(ctx context.Context, id uuid.UUID) (model.Blog, error) {
				return model.Blog{}, model.ErrBlogNotFound
			},
			expectedStatus:  http.StatusNotFound,
			expectedSuccess: false,
			expectedMessage: "Blog not found",
		},
		{
			name:   "Validation Error — ErrIdValidation from service",
			method: http.MethodGet,
			pathID: uuid.Nil.String(),
			getByIdFunc: func(ctx context.Context, id uuid.UUID) (model.Blog, error) {
				return model.Blog{}, model.ErrIdValidation
			},
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Valid ID is required",
		},

		// =====================================================================
		// INTERNAL SERVER ERROR
		// =====================================================================
		{
			name:   "Internal Server Error — database unexpected error",
			method: http.MethodGet,
			pathID: testUUID.String(),
			getByIdFunc: func(ctx context.Context, id uuid.UUID) (model.Blog, error) {
				return model.Blog{}, errors.New("database connection refused")
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedSuccess: false,
			expectedMessage: "Internal server error. Please try again later.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockBlogService{
				GetBlogByIdFunc: tt.getByIdFunc,
			}
			blogController := controller.NewBlogController(mockService)

			req := httptest.NewRequest(tt.method, "/blog/getblog/id/test", nil)
			req.SetPathValue("id", tt.pathID)
			w := httptest.NewRecorder()

			blogController.GetBlogById_Controller(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status code: got %d, want %d", w.Code, tt.expectedStatus)
			}

			if w.Header().Get("Content-Type") != "application/json" {
				t.Errorf("Content-Type: got %q, want %q", w.Header().Get("Content-Type"), "application/json")
			}

			if tt.expectData {
				var resp model.PostBlogResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if resp.Success != tt.expectedSuccess {
					t.Errorf("Success: got %v, want %v", resp.Success, tt.expectedSuccess)
				}
				if resp.Message != tt.expectedMessage {
					t.Errorf("Message: got %q, want %q", resp.Message, tt.expectedMessage)
				}
				if resp.Data.Blog_ID != testUUID {
					t.Errorf("Data Blog_ID: got %v, want %v", resp.Data.Blog_ID, testUUID)
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
