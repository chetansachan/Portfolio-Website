package test

import (
	"bytes"
	"context"
	Base64 "encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"portfolio-website/internals/controller"
	"portfolio-website/internals/model"
	"portfolio-website/internals/service"
)

func TestUploadService_FileSizeLimit(t *testing.T) {
	uploadService := service.NewUploadService(nil)

	t.Run("File exceeds 10MB limit", func(t *testing.T) {
		// 10MB + 1 byte
		oversizedBytes := make([]byte, service.MaxFileSize+1)
		b64 := Base64.StdEncoding.EncodeToString(oversizedBytes)
		payload := "data:image/png;base64," + b64

		_, err := uploadService.Upload_Service(context.Background(), payload)
		if err == nil {
			t.Fatal("expected error for oversized file, got nil")
		}
		if err != model.ErrFileTooLarge {
			t.Errorf("expected ErrFileTooLarge, got: %v", err)
		}
	})

	t.Run("Invalid MIME structure", func(t *testing.T) {
		_, err := uploadService.Upload_Service(context.Background(), "invalid-data-string")
		if err != model.ErrInvalidMimeStructure {
			t.Errorf("expected ErrInvalidMimeStructure, got: %v", err)
		}
	})

	t.Run("Disallowed extension", func(t *testing.T) {
		b64 := Base64.StdEncoding.EncodeToString([]byte("test"))
		_, err := uploadService.Upload_Service(context.Background(), "data:application/pdf;base64,"+b64)
		if err == nil || !strings.Contains(err.Error(), "is not allowed") {
			t.Errorf("expected extension not allowed error, got: %v", err)
		}
	})
}

func TestUploadController_RequestBodyLimit(t *testing.T) {
	uploadService := service.NewUploadService(nil)
	uploadController := controller.NewUploadController(uploadService)

	t.Run("Method not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/upload", nil)
		rec := httptest.NewRecorder()

		uploadController.Upload_Controller(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status 405, got %d", rec.Code)
		}
	})

	t.Run("File too large via oversized base64 payload in service", func(t *testing.T) {
		// 10MB + 1 byte
		oversizedBytes := make([]byte, service.MaxFileSize+1)
		b64 := Base64.StdEncoding.EncodeToString(oversizedBytes)
		payload := controller.Upload{
			File: "data:image/png;base64," + b64,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/upload", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		uploadController.Upload_Controller(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}

		var resp model.Response
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if resp.Success {
			t.Error("expected success to be false")
		}
		if resp.Message != "File size exceeds 10MB limit" {
			t.Errorf("expected 'File size exceeds 10MB limit', got '%s'", resp.Message)
		}
	})

	t.Run("Request body exceeds max request body size", func(t *testing.T) {
		// Valid JSON structure with >15MB content to trigger MaxBytesReader cutoff
		prefix := `{"file":"data:image/png;base64,`
		filler := strings.Repeat("A", 16*1024*1024)
		suffix := `"}`
		reader := io.MultiReader(
			strings.NewReader(prefix),
			strings.NewReader(filler),
			strings.NewReader(suffix),
		)
		req := httptest.NewRequest(http.MethodPost, "/upload", reader)
		rec := httptest.NewRecorder()

		uploadController.Upload_Controller(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}

		var resp model.Response
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if resp.Success {
			t.Error("expected success to be false")
		}
		if resp.Message != "File too large" {
			t.Errorf("expected 'File too large', got '%s'", resp.Message)
		}
	})
}
