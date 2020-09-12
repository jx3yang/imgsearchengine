package api

import (
	"encoding/json"
	"image"
	"net/http"
	"strconv"

	engine "github.com/jx3yang/imgsearchengine/src/engine"
	phash "github.com/jx3yang/imgsearchengine/src/phash"
	vptree "github.com/jx3yang/imgsearchengine/src/vptree"
)

const contentTypeKey = "Content-Type"
const defaultContentType = "application/json"

type queryResult struct {
	path     string
	distance float64
}

// EngineAPI serves the image searching engine
type EngineAPI struct {
	Tree *vptree.VPTree
}

func (service *EngineAPI) invalidTree() bool {
	return service.Tree == nil || service.Tree.Root == nil
}

// Ping will check if the engine is ready
func (service *EngineAPI) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentTypeKey, defaultContentType)

	result := make(map[string]bool)

	if service.invalidTree() {
		result["ready"] = false
	} else {
		result["ready"] = true
	}

	json.NewEncoder(w).Encode(result)
}

// KNNSearch will look for the k nearest neighbours of the given point
// where the point is expected to be the uint representation of the phash of the image
func (service *EngineAPI) KNNSearch(w http.ResponseWriter, r *http.Request) {
	k, errK := strconv.ParseUint(r.FormValue("query"), 10, 64)

	if errK != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	searchFnc := func(img image.Image) ([]map[string]interface{}, error) {
		return service.knnSearch(img, uint(k))
	}

	service.search(w, r, searchFnc)
}

func (service *EngineAPI) knnSearch(img image.Image, k uint) ([]map[string]interface{}, error) {
	searchFnc := func(queryPoint *engine.ImageInfo) (map[interface{}]float64, error) {
		return service.Tree.KNNSearch(queryPoint, uint(k))
	}

	return getResults(img, searchFnc)
}

// RangeSearch will look all the points within `threshold` distance
// of the given point
func (service *EngineAPI) RangeSearch(w http.ResponseWriter, r *http.Request) {
	threshold, errT := strconv.ParseFloat(r.FormValue("query"), 64)

	if errT != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	searchFnc := func(img image.Image) ([]map[string]interface{}, error) {
		return service.rangeSearch(img, threshold)
	}

	service.search(w, r, searchFnc)
}

func (service *EngineAPI) rangeSearch(img image.Image, threshold float64) ([]map[string]interface{}, error) {
	searchFnc := func(queryPoint *engine.ImageInfo) (map[interface{}]float64, error) {
		return service.Tree.RangeSearch(queryPoint, threshold)
	}

	return getResults(img, searchFnc)
}

func getResults(img image.Image, searchFnc func(*engine.ImageInfo) (map[interface{}]float64, error)) ([]map[string]interface{}, error) {
	hash := phash.GetPHash(img)
	queryPoint := engine.NewImageInfo(hash, "")
	searchResults, err := searchFnc(queryPoint)
	if err != nil {
		return nil, err
	}
	results := make([]map[string]interface{}, 0)
	for k, v := range searchResults {
		elem := make(map[string]interface{})
		imgInfo := k.(*engine.ImageInfo)
		imgInfoMap := map[string]interface{}{"path": imgInfo.GetPath(), "phash": imgInfo.GetPHash()}
		elem["imageInfo"] = imgInfoMap
		elem["distance"] = v
		results = append(results, elem)
	}
	return results, nil
}

func (service *EngineAPI) search(w http.ResponseWriter, r *http.Request, searchFnc func(img image.Image) ([]map[string]interface{}, error)) {
	imagePath := r.FormValue("image")

	if service.invalidTree() {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		resp, errG := http.Get(imagePath)
		if errG != nil || resp.StatusCode != http.StatusOK {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		img, _, errImg := image.Decode(resp.Body)
		if errImg != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		results, errSearch := searchFnc(img)
		if errSearch != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(results)
		}
	}
}
