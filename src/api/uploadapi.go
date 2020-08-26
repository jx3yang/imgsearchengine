package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// UploadAPI serves the image uploading service
type UploadAPI struct {
	Path string
}

// Method to complete if using external API to save the images
func (service *UploadAPI) authenticate( /**/ ) { /**/ }

// Ping will check if the service is ready
func (service *UploadAPI) Ping(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ready": true})
}

// ImageUpload handles the uploading of an image
func (service *UploadAPI) ImageUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("file")
	defer file.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// handler.Filename, handler.Size, handler.Header

	fileName := uuid.New().String() + ".jpg"
	filePath, _ := filepath.Abs(service.Path + fileName)
	out, err := os.Create(filePath)
	defer out.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, errC := io.Copy(out, file)
	if errC != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"path": filePath})
}
