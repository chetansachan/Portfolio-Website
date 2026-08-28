package controller

import (
	"encoding/json"
	"net/http"
	"portfolio-website/internals/model"
	"portfolio-website/internals/service"
)

type Upload struct {
	File string `json:"file"`
}

func Upload_Controller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.Response{
			Success: false,
			Message: "Method Not Allowed",
		})
		return
	}

	var url Upload
	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.Response{
			Success: false,
			Message: "Invalid File Type",
		})
		return
	}

	if err2 := service.Upload_Service(r.Context(), url.File); err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.Response{
			Success: false,
			Message: "File Type Not Accepted",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.Response{
		Success: true,
		Message: "File Uploaded",
	})
}

func FindFile_Controller(w http.ResponseWriter, r *http.Request) {
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

	filePath, err2 := service.FindFile_Service(r.Context(), fileName)
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
