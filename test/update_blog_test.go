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
// Tests for PATCH /blog/updateblog/{id}
// =============================================================================

func TestUpdateBlogById_Controller(t *testing.T) {
	testUUID := uuid.New()
	updatedBlog := model.Blog{
		Blog_ID:      testUUID,
		Blog_Title:   "Updated Title",
		Blog_DES:     "Updated Description",
		Blog_Slug:    "updated-slug",
		Blog_Content: "Updated content...",
		Blog_Banner:  "https://example.com/new-banner.png",
		Created_At:   time.Now(),
	}

	tests := []struct {
		name string

		// Request configuration
		method string
		pathID string
		body   string

		// Mock configuration
		updateFunc func(ctx context.Context, id uuid.UUID, title, des, slug, content, banner string) (model.Blog, error)

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
			name:   "Success — all fields updated successfully",
			method: http.MethodPatch,
			pathID: testUUID.String(),
			body: `{
				"blog_title": "Updated Title",
				"blog_des": "Updated Description",
				"blog_slug": "updated-slug",
				"blog_content": "Updated content...",
				"blog_banner": "https://example.com/new-banner.png"
			}`,
			updateFunc: func(ctx context.Context, id uuid.UUID, title, des, slug, content, banner string) (model.Blog, error) {
				return updatedBlog, nil
			},
			expectedStatus:  http.StatusOK,
			expectedSuccess: true,
			expectedMessage: "Blog Updated Successfully",
			expectData:      true,
		},
		{
			name:   "Success — partial fields provided (title and description only)",
			method: http.MethodPatch,
			pathID: testUUID.String(),
			body: `{
				"blog_title": "Updated Title Only",
				"blog_des": "Updated Description Only"
			}`,
			updateFunc: func(ctx context.Context, id uuid.UUID, title, des, slug, content, banner string) (model.Blog, error) {
				return updatedBlog, nil
			},
			expectedStatus:  http.StatusOK,
			expectedSuccess: true,
			expectedMessage: "Blog Updated Successfully",
			expectData:      true,
		},

		// =====================================================================
		// WRONG HTTP METHOD CASES
		// =====================================================================
		{
			name:            "Method Not Allowed — GET instead of PATCH",
			method:          http.MethodGet,
			pathID:          testUUID.String(),
			body:            "",
			updateFunc:      nil,
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedSuccess: false,
			expectedMessage: "Method not allowed",
		},
		{
			name:            "Method Not Allowed — POST instead of PATCH",
			method:          http.MethodPost,
			pathID:          testUUID.String(),
			body:            "",
			updateFunc:      nil,
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedSuccess: false,
			expectedMessage: "Method not allowed",
		},
		{
			name:            "Method Not Allowed — DELETE instead of PATCH",
			method:          http.MethodDelete,
			pathID:          testUUID.String(),
			body:            "",
			updateFunc:      nil,
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedSuccess: false,
			expectedMessage: "Method not allowed",
		},

		// =====================================================================
		// INVALID ID FORMAT CASES
		// =====================================================================
		{
			name:            "Bad Request — non-UUID ID",
			method:          http.MethodPatch,
			pathID:          "not-a-uuid",
			body:            `{"blog_title": "Test"}`,
			updateFunc:      nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Invalid ID format",
		},
		{
			name:            "Bad Request — empty ID",
			method:          http.MethodPatch,
			pathID:          "",
			body:            `{"blog_title": "Test"}`,
			updateFunc:      nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Invalid ID format",
		},

		// =====================================================================
		// NOT FOUND / VALIDATION / SERVICE ERRORS
		// =====================================================================
		{
			name:   "Not Found — blog does not exist to update",
			method: http.MethodPatch,
			pathID: uuid.New().String(),
			body:   `{"blog_title": "New Title"}`,
			updateFunc: func(ctx context.Context, id uuid.UUID, title, des, slug, content, banner string) (model.Blog, error) {
				return model.Blog{}, model.ErrBlogNotFound
			},
			expectedStatus:  http.StatusNotFound,
			expectedSuccess: false,
			expectedMessage: "Blog Not Found",
		},
		{
			name:   "Validation Error — ErrIdValidation from service",
			method: http.MethodPatch,
			pathID: uuid.Nil.String(),
			body:   `{"blog_title": "New Title"}`,
			updateFunc: func(ctx context.Context, id uuid.UUID, title, des, slug, content, banner string) (model.Blog, error) {
				return model.Blog{}, model.ErrIdValidation
			},
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Id field is required",
		},
		{
			name:   "Internal Server Error — database update failed",
			method: http.MethodPatch,
			pathID: testUUID.String(),
			body:   `{"blog_title": "New Title"}`,
			updateFunc: func(ctx context.Context, id uuid.UUID, title, des, slug, content, banner string) (model.Blog, error) {
				return model.Blog{}, errors.New("db write failed")
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedSuccess: false,
			expectedMessage: "Internal server error. Please try again later.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockBlogService{
				UpdateBlogFunc: func(ctx context.Context, id uuid.UUID, title string, des string, slug string, content string, banner string, featured bool) (model.Blog, error) {
					if tt.updateFunc == nil {
						return model.Blog{}, nil
					}
					return tt.updateFunc(ctx, id, title, des, slug, content, banner)
				},
			}
			blogController := controller.NewBlogController(mockService)

			req := httptest.NewRequest(tt.method, "/blog/updateblog/test", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			req.SetPathValue("id", tt.pathID)
			w := httptest.NewRecorder()

			blogController.UpdateBlogById_Controller(w, req)

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
