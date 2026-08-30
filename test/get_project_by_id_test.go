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
// Tests for GET /project/getproject/id/{id}
// =============================================================================

func TestGetProjectById_Controller(t *testing.T) {
	testUUID := uuid.New()
	sampleProject := model.Project{
		Project_ID:      testUUID,
		Project_Name:    "Awesome Project",
		Project_DES:     "Project description",
		Project_Banner:  "https://example.com/banner.png",
		Project_Content: "Detailed project documentation",
		Project_Slug:    "awesome-project",
		Is_Featured:     true,
		Created_At:      time.Now(),
	}

	tests := []struct {
		name string

		// Request configuration
		method string
		pathID string

		// Mock configuration
		getByIdFunc func(ctx context.Context, id uuid.UUID) (model.Project, error)

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
			name:   "Success — project found by ID",
			method: http.MethodGet,
			pathID: testUUID.String(),
			getByIdFunc: func(ctx context.Context, id uuid.UUID) (model.Project, error) {
				return sampleProject, nil
			},
			expectedStatus:  http.StatusOK,
			expectedSuccess: true,
			expectedMessage: "Project Found Successfully",
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
			expectedMessage: "Method Not Allowed",
		},
		{
			name:            "Method Not Allowed — DELETE instead of GET",
			method:          http.MethodDelete,
			pathID:          testUUID.String(),
			getByIdFunc:     nil,
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedSuccess: false,
			expectedMessage: "Method Not Allowed",
		},

		// =====================================================================
		// INVALID ID FORMAT CASES
		// =====================================================================
		{
			name:            "Bad Request — non-UUID format",
			method:          http.MethodGet,
			pathID:          "invalid-id",
			getByIdFunc:     nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Invalid Id Format",
		},
		{
			name:            "Bad Request — empty ID string",
			method:          http.MethodGet,
			pathID:          "",
			getByIdFunc:     nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Invalid Id Format",
		},

		// =====================================================================
		// NOT FOUND / VALIDATION ERRORS
		// =====================================================================
		{
			name:   "Bad Request / Not Found — project does not exist",
			method: http.MethodGet,
			pathID: uuid.New().String(),
			getByIdFunc: func(ctx context.Context, id uuid.UUID) (model.Project, error) {
				return model.Project{}, model.ErrProjectNotFound
			},
			expectedStatus:  http.StatusNotFound,
			expectedSuccess: false,
			expectedMessage: "Project Not Found",
		},
		{
			name:   "Validation Error — ErrIdValidation",
			method: http.MethodGet,
			pathID: uuid.Nil.String(),
			getByIdFunc: func(ctx context.Context, id uuid.UUID) (model.Project, error) {
				return model.Project{}, model.ErrIdValidation
			},
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Id is Required",
		},

		// =====================================================================
		// INTERNAL SERVER ERROR
		// =====================================================================
		{
			name:   "Internal Server Error — database unexpected error",
			method: http.MethodGet,
			pathID: testUUID.String(),
			getByIdFunc: func(ctx context.Context, id uuid.UUID) (model.Project, error) {
				return model.Project{}, errors.New("db error")
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedSuccess: false,
			expectedMessage: "Internal Server Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockProjectService{
				GetProjectByIdFunc: tt.getByIdFunc,
			}
			projectController := controller.NewProjectController(mockService)

			req := httptest.NewRequest(tt.method, "/project/getproject/id/test", nil)
			req.SetPathValue("id", tt.pathID)
			w := httptest.NewRecorder()

			projectController.GetProjectById_Controller(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status code: got %d, want %d", w.Code, tt.expectedStatus)
			}

			if w.Header().Get("Content-Type") != "application/json" {
				t.Errorf("Content-Type: got %q, want %q", w.Header().Get("Content-Type"), "application/json")
			}

			if tt.expectData {
				var resp model.PostProjectResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if resp.Success != tt.expectedSuccess {
					t.Errorf("Success: got %v, want %v", resp.Success, tt.expectedSuccess)
				}
				if resp.Message != tt.expectedMessage {
					t.Errorf("Message: got %q, want %q", resp.Message, tt.expectedMessage)
				}
				if resp.Data.Project_ID != testUUID {
					t.Errorf("Data Project_ID: got %v, want %v", resp.Data.Project_ID, testUUID)
				}
			} else {
				var resp model.ProjectResponse
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
