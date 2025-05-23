package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"mqtt-streaming-server/domain"
	"mqtt-streaming-server/repository"
	"mqtt-streaming-server/utils"
)

type PhotoController struct {
	PhotoRepository domain.PhotoRepository
}

func InitPhotoRoutes(db *mongo.Database, mux *http.ServeMux) {
	photoController := &PhotoController{
		PhotoRepository: repository.NewPhotoRepository(db),
	}

	mux.Handle("/photos", withAuth(http.HandlerFunc(photoController.GetPhotos)))
}

func (ctlr PhotoController) GetPhotos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	text := r.URL.Query().Get("text")
	deviceID := r.URL.Query().Get("device_id")

	if start == "" {
		start = strconv.FormatInt(time.Now().Add(-24*time.Hour).UTC().Unix(), 10)
	}

	if end == "" {
		end = strconv.FormatInt(time.Now().UTC().Unix(), 10)
	}

	startInt, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		http.Error(w, "Invalid start timestamp "+err.Error(), http.StatusBadRequest)
		return
	}

	endInt, err := strconv.ParseInt(end, 10, 64)
	if err != nil {
		http.Error(w, "Invalid end timestamp "+err.Error(), http.StatusBadRequest)
		return
	}

	filters := map[string]any{
		"timestamp": map[string]any{
			"$gte": time.Unix(startInt, 0),
			"$lte": time.Unix(endInt, 0),
		},
	}

	if text != "" {
		filters["text"] = map[string]any{
			"$regex":   text,
			"$options": "i",
		}
	}

	if deviceID != "" {
		filters["device_id"] = deviceID
	}

	photos, err := ctlr.PhotoRepository.GetPhotos(ctx, filters)
	if err != nil {
		fmt.Println("Error fetching photos:", err)
		http.Error(w, "Failed to fetch photos: ", http.StatusInternalServerError)
		return
	}

	for _, photo := range photos {
		presignedURL, err := utils.GetPresignedURL(ctx, fmt.Sprintf("photos/%d.%s", photo.Timestamp.Unix(), photo.ImageType))
		if err != nil {
			http.Error(w, "Failed to get presigned URL", http.StatusInternalServerError)
			return
		}
		photo.PresignedURL = presignedURL
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(photos)
}
