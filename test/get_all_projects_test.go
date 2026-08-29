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
// Tests for GET /project/getallprojects
// =============================================================================

func TestGetAllProjects_Controller(t *testing.T) {
	sampleProjects := []model.Project{
		{
			Project_ID:      uuid.New(),
			Project_Name:    "Project Alpha",
			Project_DES:     "First Project",
			Project_Banner:  "https://example.com/alpha.png",
			Project_Content: "Content Alpha",
			Project_Slug:    "project-alpha",
			Is_Featured:     true,
			Created_At:      time.Now(),
		},
		{
			Project_ID:      uuid.New(),
			Project_Name:    "Project Beta",
			Project_DES:     "Second Project",
			Project_Banner:  "https://example.com/beta.png",
			Project_Content: "Content Beta",
			Project_Slug:    "project-beta",
			Is_Featured:     false,
			Created_At:      time.Now().Add(-time.Hour),
		},
	}

	tests := []struct {
		name string

		// Request configuration
		method      string
		queryParams string

		// Mock configuration
		getAllFunc func(ctx context.Context, page int, offset int, featured bool) ([]model.Project, int, error)

		// Expected response
		expectedStatus  int
		expectedSuccess bool
		expectDataCount int
		expectTotal     int
	}{
		// =====================================================================
		// HAPPY PATH
		// =====================================================================
		{
			name:        "Success — default query params",
			method:      http.MethodGet,
			queryParams: "",
			getAllFunc: func(ctx context.Context, page int, offset int, featured bool) ([]model.Project, int, error) {
				return sampleProjects, 2, nil
			},
			expectedStatus:  http.StatusOK,
			expectedSuccess: true,
			expectDataCount: 2,
			expectTotal:     2,
		},
		{
			name:        "Success — valid custom pagination",
			method:      http.MethodGet,
			queryParams: "?page=2&size=5",
			getAllFunc: func(ctx context.Context, page int, offset int, featured bool) ([]model.Project, int, error) {
				return sampleProjects[:1], 10, nil
			},
			expectedStatus:  http.StatusOK,
			expectedSuccess: true,
			expectDataCount: 1,
			expectTotal:     10,
		},
		{
			name:        "Success — size capped at 100",
			method:      http.MethodGet,
			queryParams: "?page=1&size=200",
			getAllFunc: func(ctx context.Context, page int, offset int, featured bool) ([]model.Project, int, error) {
				return sampleProjects, 2, nil
			},
			expectedStatus:  http.StatusOK,
			expectedSuccess: true,
			expectDataCount: 2,
			expectTotal:     2,
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
		},

		// =====================================================================
		// INVALID PARAMETERS
		// =====================================================================
		{
			name:            "Bad Request — non-numeric page",
			method:          http.MethodGet,
			queryParams:     "?page=abc",
			getAllFunc:      nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
		},
		{
			name:            "Bad Request — zero page",
			method:          http.MethodGet,
			queryParams:     "?page=0",
			getAllFunc:      nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
		},
		{
			name:            "Bad Request — negative page",
			method:          http.MethodGet,
			queryParams:     "?page=-1",
			getAllFunc:      nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
		},
		{
			name:            "Bad Request — non-numeric size",
			method:          http.MethodGet,
			queryParams:     "?size=xyz",
			getAllFunc:      nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
		},
		{
			name:            "Bad Request — zero size",
			method:          http.MethodGet,
			queryParams:     "?size=0",
			getAllFunc:      nil,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
		},

		// =====================================================================
		// NOT FOUND / SERVICE ERRORS
		// =====================================================================
		{
			name:        "Not Found — service returns ErrProjectNotFound",
			method:      http.MethodGet,
			queryParams: "",
			getAllFunc: func(ctx context.Context, page int, offset int, featured bool) ([]model.Project, int, error) {
				return nil, 0, model.ErrProjectNotFound
			},
			expectedStatus:  http.StatusNotFound,
			expectedSuccess: false,
		},
		{
			name:        "Internal Server Error — database error",
			method:      http.MethodGet,
			queryParams: "",
			getAllFunc: func(ctx context.Context, page int, offset int, featured bool) ([]model.Project, int, error) {
				return nil, 0, errors.New("database unavailable")
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockProjectService{
				GetAllProjectsFunc: tt.getAllFunc,
			}
			projectController := controller.NewProjectController(mockService)

			req := httptest.NewRequest(tt.method, "/project/getallprojects"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			projectController.GetAllProjects_Controller(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status code: got %d, want %d", w.Code, tt.expectedStatus)
			}

			if w.Header().Get("Content-Type") != "application/json" {
				t.Errorf("Content-Type: got %q, want %q", w.Header().Get("Content-Type"), "application/json")
			}

			if tt.expectedStatus == http.StatusOK {
				var resp model.PaginatedProjectResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if resp.Success != tt.expectedSuccess {
					t.Errorf("Success: got %v, want %v", resp.Success, tt.expectedSuccess)
				}
				if len(resp.Data) != tt.expectDataCount {
					t.Errorf("Data count: got %d, want %d", len(resp.Data), tt.expectDataCount)
				}
				if resp.Meta.TotalProjects != tt.expectTotal {
					t.Errorf("Meta TotalProjects: got %d, want %d", resp.Meta.TotalProjects, tt.expectTotal)
				}
			}
		})
	}
}
