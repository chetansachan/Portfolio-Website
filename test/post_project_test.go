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
// Tests for POST /project/createproject
// =============================================================================

func TestCreateProject_Controller(t *testing.T) {
	sampleProject := model.Project{
		Project_ID:      uuid.New(),
		Project_Name:    "Portfolio Website",
		Project_DES:     "A personal portfolio with Go and Next.js",
		Project_Banner:  "https://example.com/project-banner.png",
		Project_Content: "Full project details...",
		Project_Slug:    "portfolio-website",
		Is_Featured:     true,
		Created_At:      time.Now(),
	}

	tests := []struct {
		name string

		// Request configuration
		method string
		body   string

		// Mock configuration
		createFunc func(ctx context.Context, name, des, banner, content, slug string, featured bool) (model.Project, error)

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
			name:   "Success — all fields provided, project created",
			method: http.MethodPost,
			body: `{
				"project_name": "Portfolio Website",
				"project_des": "A personal portfolio with Go and Next.js",
				"project_banner": "https://example.com/project-banner.png",
				"project_content": "Full project details...",
				"project_slug": "portfolio-website",
				"is_featured": true
			}`,
			createFunc: func(ctx context.Context, name, des, banner, content, slug string, featured bool) (model.Project, error) {
				return sampleProject, nil
			},
			expectedStatus:  http.StatusCreated,
			expectedSuccess: true,
			expectedMessage: "Project created successfully",
			expectData:      true,
		},

		// =====================================================================
		// WRONG HTTP METHOD CASES
		// =====================================================================
		{
			name:            "Method Not Allowed — GET instead of POST",
			method:          http.MethodGet,
			body:            "",
			createFunc:      nil,
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedSuccess: false,
			expectedMessage: "Method Not Allowed",
		},
		{
			name:            "Method Not Allowed — DELETE instead of POST",
			method:          http.MethodDelete,
			body:            "",
			createFunc:      nil,
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedSuccess: false,
			expectedMessage: "Method Not Allowed",
		},
		{
			name:            "Method Not Allowed — PATCH instead of POST",
			method:          http.MethodPatch,
			body:            "",
			createFunc:      nil,
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedSuccess: false,
			expectedMessage: "Method Not Allowed",
		},

		// =====================================================================
		// MALFORMED / INVALID JSON BODY
		// =====================================================================
		{
			name:            "Bad Request — invalid JSON body",
			method:          http.MethodPost,
			body:            `not valid json`,
			createFunc:      nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Invalid request body",
		},
		{
			name:            "Bad Request — empty body",
			method:          http.MethodPost,
			body:            "",
			createFunc:      nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "Invalid request body",
		},

		// =====================================================================
		// VALIDATION FAILURES
		// =====================================================================
		{
			name:   "Validation Error — empty required fields",
			method: http.MethodPost,
			body: `{
				"project_name": "",
				"project_des": "",
				"project_slug": ""
			}`,
			createFunc: func(ctx context.Context, name, des, banner, content, slug string, featured bool) (model.Project, error) {
				return model.Project{}, model.ErrValidation
			},
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedMessage: "All fields are required",
		},
		{
			name:   "Conflict — project already exists",
			method: http.MethodPost,
			body: `{
				"project_name": "Portfolio Website",
				"project_des": "Description",
				"project_banner": "https://example.com/b.png",
				"project_content": "Content",
				"project_slug": "portfolio-website",
				"is_featured": true
			}`,
			createFunc: func(ctx context.Context, name, des, banner, content, slug string, featured bool) (model.Project, error) {
				return model.Project{}, model.ErrProjectAlreadyExist
			},
			expectedStatus:  http.StatusConflict,
			expectedSuccess: false,
			expectedMessage: "Project Already Exist",
		},

		// =====================================================================
		// INTERNAL SERVER ERROR
		// =====================================================================
		{
			name:   "Internal Server Error — database failure",
			method: http.MethodPost,
			body: `{
				"project_name": "Portfolio Website",
				"project_des": "Description",
				"project_banner": "https://example.com/b.png",
				"project_content": "Content",
				"project_slug": "portfolio-website",
				"is_featured": true
			}`,
			createFunc: func(ctx context.Context, name, des, banner, content, slug string, featured bool) (model.Project, error) {
				return model.Project{}, errors.New("db error")
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedSuccess: false,
			expectedMessage: "Internal server error. Please try again later.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockProjectService{
				CreateProjectFunc: tt.createFunc,
			}
			projectController := controller.NewProjectController(mockService)

			req := httptest.NewRequest(tt.method, "/project/createproject", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			projectController.CreateProject_Controller(w, req)

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
				if resp.Data.Project_ID == uuid.Nil {
					t.Error("Expected valid Project_ID, got nil UUID")
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
