package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"portfolio-website/internals/model"
	"portfolio-website/internals/service"
)

type Upload struct {
	File string `json:"file"`
}

type UploadResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	File    string `json:"file"`
}

// maxRequestBodySize allows up to 30 MB for the JSON request body,
// which accommodates a 20 MB file encoded as Base64 JSON (~26.7 MB + metadata).
const maxRequestBodySize = 30 * 1024 * 1024

type UploadController struct {
	UploadService *service.UploadService
}

func NewUploadController(uploadService *service.UploadService) *UploadController {
	return &UploadController{UploadService: uploadService}
}

func (c *UploadController) Upload_Controller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.Response{
			Success: false,
			Message: "Method Not Allowed",
		})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)

	var File Upload
	if err := json.NewDecoder(r.Body).Decode(&File); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "File too large",
			})
			return
		}
		json.NewEncoder(w).Encode(model.Response{
			Success: false,
			Message: "Invalid Request Body",
		})
		return
	}

	file, err := c.UploadService.Upload_Service(r.Context(), File.File)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidMimeStructure),
			errors.Is(err, model.ErrInvalidMimeDetails),
			errors.Is(err, model.ErrExtensionNotAllowed),
			errors.Is(err, model.ErrDecodingBase64),
			errors.Is(err, model.ErrFileTooLarge):
			log.Printf("[WARN] File upload validation failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: err.Error(),
			})
		default:
			log.Printf("[ERROR] File upload failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "Internal Server Error",
			})
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UploadResponse{
		Success: true,
		Message: "File Uploaded",
		File:    file,
	})
}
