package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"portfolio-website/internals/model"
	"portfolio-website/internals/service"
)

type Upload struct {
	File string `json:"file"`
}

type UploadResponse struct {
	Success bool
	Message string
	Url     string
}

const maxFileSize = 10 * 1024 * 1024

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

	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)

	var file Upload
	if err := json.NewDecoder(r.Body).Decode(&file); err != nil {
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

	url, err := c.UploadService.Upload_Service(r.Context(), file.File)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.Response{
			Success: false,
			Message: "File Type Not Accepted",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UploadResponse{
		Success: true,
		Message: "File Uploaded",
		Url:     url,
	})
}

func (c *UploadController) FindFile_Controller(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.Response{
			Success: false,
			Message: "Method Not Allowed",
		})
		return
	}

	fileName := r.PathValue("file")

	filePath, err2 := c.UploadService.FindFile_Service(r.Context(), fileName)
	if err2 != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.Response{
			Success: false,
			Message: "File Not Exist",
		})
		return
	}

	http.ServeFile(w, r, filePath)
}

func (c *UploadController) DeleteFile_Controller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.Response{
			Success: false,
			Message: "Method Not Allowed",
		})
		return
	}

	fileName := r.PathValue("file")
	err := c.UploadService.DeleteFile_Service(r.Context(), fileName)
	if err != nil {
		if errors.Is(err, model.ErrFileNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.Response{
				Success: false,
				Message: "File Not Exist",
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.Response{
			Success: false,
			Message: "Internal server error. Please try again later.",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.Response{
		Success: true,
		Message: "File Deleted Successfully",
	})
}
