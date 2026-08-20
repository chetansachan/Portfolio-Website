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
// Tests for DELETE /project/deleteproject/{id}
// =============================================================================

func TestDeleteProjectById_Controller(t *testing.T) {
	validUUID := uuid.New().String()

	tests := []struct {
		name string

		// Request configuration
		method string
		pathID string

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
			name:   "Success — project exists and is deleted",
			method: http.MethodDelete,
			pathID: validUUID,
			deleteFunc: func(ctx context.Context, id uuid.UUID) error {
				return nil
			},
			expectedStatus:  http.StatusOK,
			expectedSuccess: true,
			expectedMessage: "Project Deleted Successfully",
		},

		// =====================================================================
		// NOT FOUND
		// =====================================================================
		{
			name:   "Not Found — project does not exist",
			method: http.MethodDelete,
			pathID: uuid.New().String(),
			deleteFunc: func(ctx context.Context, id uuid.UUID) error {
				return model.ErrProjectNotFound
			},
			expectedStatus:  http.StatusNotFound,
			expectedSuccess: false,
			expectedMessage: "Project Not Found",
		},

		// =====================================================================
		// INVALID ID FORMAT CASES
		// =====================================================================
		{
			name:            "Bad Request — non-UUID format",
			method:          http.MethodDelete,
			pathID:          "not-a-valid-uuid",
			deleteFunc:      nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Invalid ID format",
		},
		{
			name:            "Bad Request — empty ID",
			method:          http.MethodDelete,
			pathID:          "",
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

		// =====================================================================
		// INTERNAL SERVER ERROR
		// =====================================================================
		{
			name:   "Internal Server Error — database error",
			method: http.MethodDelete,
			pathID: validUUID,
			deleteFunc: func(ctx context.Context, id uuid.UUID) error {
				return errors.New("db error")
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedSuccess: false,
			expectedMessage: "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockProjectService{
				DeleteProjectFunc: tt.deleteFunc,
			}
			projectController := controller.NewProjectController(mockService)

			req := httptest.NewRequest(tt.method, "/project/deleteproject/test", nil)
			req.SetPathValue("id", tt.pathID)
			w := httptest.NewRecorder()

			projectController.DeleteProjectById_Controller(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status code: got %d, want %d", w.Code, tt.expectedStatus)
			}

			if w.Header().Get("Content-Type") != "application/json" {
				t.Errorf("Content-Type: got %q, want %q", w.Header().Get("Content-Type"), "application/json")
			}

			var resp model.Response
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if resp.Success != tt.expectedSuccess {
				t.Errorf("Success: got %v, want %v", resp.Success, tt.expectedSuccess)
			}

			if resp.Message != tt.expectedMessage {
				t.Errorf("Message: got %q, want %q", resp.Message, tt.expectedMessage)
			}
		})
	}
}
