package test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"portfolio-website/internals/controller"
	"portfolio-website/internals/model"

	"github.com/google/uuid"
)

// =============================================================================
// Tests for DELETE /blog/deleteblog/{id}
// =============================================================================

func TestDeleteBlogById_Controller(t *testing.T) {
	validUUID := uuid.New().String()

	tests := []struct {
		name string // descriptive test case name

		// Request configuration
		method string // HTTP method to use
		pathID string // value for {id} path parameter

		// Mock configuration
		deleteFunc func(ctx context.Context, id uuid.UUID) error

		// Expected response
		expectedStatus  int
		expectedSuccess bool
		expectedMessage string
	}{
		// =====================================================================
		// HAPPY PATH
		// =====================================================================
		{
			name:   "Success — blog exists and is deleted",
			method: http.MethodDelete,
			pathID: validUUID,
			deleteFunc: func(ctx context.Context, id uuid.UUID) error {
				return nil // simulate successful deletion
			},
			expectedStatus:  http.StatusOK,
			expectedSuccess: true,
			expectedMessage: "Blog Deleted Successfully",
		},

		// =====================================================================
		// BLOG NOT FOUND CASES
		// =====================================================================
		{
			name:   "Not Found — blog does not exist (never created)",
			method: http.MethodDelete,
			pathID: uuid.New().String(), // random UUID that doesn't exist
			deleteFunc: func(ctx context.Context, id uuid.UUID) error {
				return model.ErrBlogNotFound
			},
			expectedStatus:  http.StatusNotFound,
			expectedSuccess: false,
			expectedMessage: "Blog Not Found",
		},
		{
			name:   "Not Found — blog was already deleted (idempotency check)",
			method: http.MethodDelete,
			pathID: validUUID, // same UUID as the happy path
			deleteFunc: func(ctx context.Context, id uuid.UUID) error {
				return model.ErrBlogNotFound // second delete should fail
			},
			expectedStatus:  http.StatusNotFound,
			expectedSuccess: false,
			expectedMessage: "Blog Not Found",
		},

		// =====================================================================
		// INVALID UUID FORMAT CASES
		// =====================================================================
		{
			name:            "Bad Request — completely invalid UUID string",
			method:          http.MethodDelete,
			pathID:          "not-a-valid-uuid",
			deleteFunc:      nil, // service should never be called
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Invalid ID format",
		},
		{
			name:            "Bad Request — empty ID parameter",
			method:          http.MethodDelete,
			pathID:          "",
			deleteFunc:      nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Invalid ID format",
		},
		{
			name:            "Bad Request — UUID-like but wrong length (truncated)",
			method:          http.MethodDelete,
			pathID:          "755cc3fb-ba02-4e69-9cdf", // truncated UUID
			deleteFunc:      nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Invalid ID format",
		},
		{
			name:            "Bad Request — SQL injection attempt in ID",
			method:          http.MethodDelete,
			pathID:          "1 OR 1=1",
			deleteFunc:      nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Invalid ID format",
		},
		{
			name:            "Bad Request — numeric ID instead of UUID",
			method:          http.MethodDelete,
			pathID:          "12345",
			deleteFunc:      nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Invalid ID format",
		},
		{
			name:            "Bad Request — whitespace-only ID",
			method:          http.MethodDelete,
			pathID:          "   ",
			deleteFunc:      nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Invalid ID format",
		},

		// =====================================================================
		// WRONG HTTP METHOD CASES
		// =====================================================================
		{
			name:            "Method Not Allowed — GET instead of DELETE",
			method:          http.MethodGet,
			pathID:          validUUID,
			deleteFunc:      nil,
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedSuccess: false,
			expectedMessage: "Method not allowed",
		},
		{
			name:            "Method Not Allowed — POST instead of DELETE",
			method:          http.MethodPost,
			pathID:          validUUID,
			deleteFunc:      nil,
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedSuccess: false,
			expectedMessage: "Method not allowed",
		},
		{
			name:            "Method Not Allowed — PUT instead of DELETE",
			method:          http.MethodPut,
			pathID:          validUUID,
			deleteFunc:      nil,
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedSuccess: false,
			expectedMessage: "Method not allowed",
		},
		{
			name:            "Method Not Allowed — PATCH instead of DELETE",
			method:          http.MethodPatch,
			pathID:          validUUID,
			deleteFunc:      nil,
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedSuccess: false,
			expectedMessage: "Method not allowed",
		},

		// =====================================================================
		// INTERNAL SERVER ERROR (unexpected failures)
		// =====================================================================
		{
			name:   "Internal Server Error — database connection dropped",
			method: http.MethodDelete,
			pathID: validUUID,
			deleteFunc: func(ctx context.Context, id uuid.UUID) error {
				return errors.New("connection refused")
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedSuccess: false,
			expectedMessage: "Internal server error. Please try again later.",
		},
		{
			name:   "Internal Server Error — unexpected unknown error",
			method: http.MethodDelete,
			pathID: validUUID,
			deleteFunc: func(ctx context.Context, id uuid.UUID) error {
				return errors.New("something completely unexpected")
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedSuccess: false,
			expectedMessage: "Internal server error. Please try again later.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// --- Setup mock ---
			mockService := &MockBlogService{
				DeleteFunc: tt.deleteFunc,
			}
			blogController := controller.NewBlogController(mockService)

			// --- Build request ---
			// Use a clean URL path; the pathID is injected via SetPathValue
			// (which is exactly how Go's ServeMux delivers it in production).
			req := httptest.NewRequest(tt.method, "/blog/deleteblog/test-id", nil)
			req.SetPathValue("id", tt.pathID)

			// --- Record response ---
			w := httptest.NewRecorder()
			blogController.DeleteBlogById_Controller(w, req)

			// --- Assert status code ---
			if w.Code != tt.expectedStatus {
				t.Errorf("Status code: got %d, want %d", w.Code, tt.expectedStatus)
			}

			// --- Assert response body ---
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

			// --- Assert Content-Type header ---
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Content-Type: got %q, want %q", contentType, "application/json")
			}
		})
	}
}
