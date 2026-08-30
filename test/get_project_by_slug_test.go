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
// Tests for GET /project/getproject/slug/{slug}
// =============================================================================

func TestGetProjectBySlug_Controller(t *testing.T) {
	sampleProject := model.Project{
		Project_ID:      uuid.New(),
		Project_Name:    "Golang Microservices",
		Project_DES:     "High-throughput microservices architecture",
		Project_Banner:  "https://example.com/microservices.png",
		Project_Content: "Detailed architecture diagram and implementation...",
		Project_Slug:    "golang-microservices",
		Is_Featured:     true,
		Created_At:      time.Now(),
	}

	tests := []struct {
		name string

		// Request configuration
		method   string
		pathSlug string

		// Mock configuration
		getBySlugFunc func(ctx context.Context, slug string) (model.Project, error)

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
			name:     "Success — project found by slug",
			method:   http.MethodGet,
			pathSlug: "golang-microservices",
			getBySlugFunc: func(ctx context.Context, slug string) (model.Project, error) {
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
			pathSlug:        "golang-microservices",
			getBySlugFunc:   nil,
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedSuccess: false,
			expectedMessage: "Method Not Allowed",
		},
		{
			name:            "Method Not Allowed — DELETE instead of GET",
			method:          http.MethodDelete,
			pathSlug:        "golang-microservices",
			getBySlugFunc:   nil,
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedSuccess: false,
			expectedMessage: "Method Not Allowed",
		},

		// =====================================================================
		// VALIDATION / NOT FOUND CASES
		// =====================================================================
		{
			name:     "Validation Error — empty slug",
			method:   http.MethodGet,
			pathSlug: "",
			getBySlugFunc: func(ctx context.Context, slug string) (model.Project, error) {
				return model.Project{}, model.ErrSlugValidation
			},
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Slug is Required",
		},
		{
			name:     "Not Found — project slug does not exist",
			method:   http.MethodGet,
			pathSlug: "non-existent-project-slug",
			getBySlugFunc: func(ctx context.Context, slug string) (model.Project, error) {
				return model.Project{}, model.ErrProjectNotFound
			},
			expectedStatus:  http.StatusNotFound,
			expectedSuccess: false,
			expectedMessage: "Project Not Found",
		},

		// =====================================================================
		// INTERNAL SERVER ERROR
		// =====================================================================
		{
			name:     "Internal Server Error — database error",
			method:   http.MethodGet,
			pathSlug: "golang-microservices",
			getBySlugFunc: func(ctx context.Context, slug string) (model.Project, error) {
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
				GetProjectBySlugFunc: tt.getBySlugFunc,
			}
			projectController := controller.NewProjectController(mockService)

			req := httptest.NewRequest(tt.method, "/project/getproject/slug/test", nil)
			req.SetPathValue("slug", tt.pathSlug)
			w := httptest.NewRecorder()

			projectController.GetProjectBySlug_Controller(w, req)

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
				if resp.Data.Project_Slug != tt.pathSlug {
					t.Errorf("Data Project_Slug: got %q, want %q", resp.Data.Project_Slug, tt.pathSlug)
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
