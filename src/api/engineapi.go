package api

import (
	"encoding/json"
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
	point, errP := strconv.ParseUint(r.FormValue("point"), 10, 64)
	k, errK := strconv.ParseUint(r.FormValue("k"), 10, 64)

	if errP != nil || errK != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if service.invalidTree() {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		hash := phash.PHash(point)
		queryPoint := engine.NewImageInfo(hash, "")
		knnMap, errKnn := service.Tree.KNNSearch(queryPoint, uint(k))
		if errKnn != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		results := make([]map[string]interface{}, 0)
		for k, v := range knnMap {
			elem := make(map[string]interface{})
			imgInfo := k.(*engine.ImageInfo)
			imgInfoMap := map[string]interface{}{"path": imgInfo.GetPath(), "phash": imgInfo.GetPHash()}
			elem["imageInfo"] = imgInfoMap
			elem["distance"] = v
			results = append(results, elem)
		}
		json.NewEncoder(w).Encode(results)
	}
}

// RangeSearch will look all the points within `threshold` distance
// of the given point
func (service *EngineAPI) RangeSearch(w http.ResponseWriter, r *http.Request) {
	point, errP := strconv.ParseUint(r.FormValue("point"), 10, 64)
	threshold, errT := strconv.ParseFloat(r.FormValue("threshold"), 64)

	if errP != nil || errT != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if service.invalidTree() {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		hash := phash.PHash(point)
		queryPoint := engine.NewImageInfo(hash, "")
		rangeMap, errR := service.Tree.RangeSearch(queryPoint, threshold)
		if errR != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		results := make([]map[string]interface{}, 0)
		for k, v := range rangeMap {
			elem := make(map[string]interface{})
			imgInfo := k.(*engine.ImageInfo)
			imgInfoMap := map[string]interface{}{"path": imgInfo.GetPath(), "phash": imgInfo.GetPHash()}
			elem["imageInfo"] = imgInfoMap
			elem["distance"] = v
			results = append(results, elem)
		}
		json.NewEncoder(w).Encode(results)
	}
}
