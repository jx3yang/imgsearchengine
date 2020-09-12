package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/jx3yang/imgsearchengine/src/api"
	"github.com/jx3yang/imgsearchengine/src/engine"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// internal paths
const imagePath = "images/"
const tempImagesPath = imagePath + "temp/"

// external paths
const port = "8080"
const devAddress = "http://localhost:" + port
const pathPrefix = "/images/"

func ping(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ready": true})
}

// imageUpload handles the uploading of an image
func imageUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("file")
	defer file.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fileName := uuid.New().String() + ".jpg"
	internalPath := tempImagesPath + fileName
	out, err := os.Create(internalPath)
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
	json.NewEncoder(w).Encode(map[string]string{"path": devAddress + "/" + internalPath})
}

func main() {
	// tree, err := engine.LoadFromCSV("load_file.csv", '\t')
	tree, err := engine.LoadFromCSVPHash("load_file_phash.csv", '\t')

	if err != nil {
		log.Fatal(err)
	}

	engineService := api.EngineAPI{Tree: tree}

	router := mux.NewRouter().StrictSlash(true)

	// Image search service
	router.HandleFunc("/knn", engineService.KNNSearch).
		Methods("POST")

	router.HandleFunc("/rangesearch", engineService.RangeSearch).
		Methods("POST")

	router.HandleFunc("/ping-engine", engineService.Ping).
		Methods("GET")

	// Dummy file server
	fs := http.FileServer(http.Dir(imagePath))
	router.PathPrefix(pathPrefix).Handler(http.StripPrefix(pathPrefix, fs))

	router.HandleFunc("/ping-image", ping).
		Methods("GET")

	router.HandleFunc("/image-upload", imageUpload).
		Methods("POST")

	log.Fatal(http.ListenAndServe(":"+port, router))
}
