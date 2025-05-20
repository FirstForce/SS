package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"mqtt-streaming-server/utils"
)

type Photo struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Timestamp    time.Time          `json:"timestamp" bson:"timestamp"`
	ImageType    string             `json:"image_type" bson:"image_type"`
	PresignedURL string             `json:"presigned_url" bson:",omitempty"`
}

type PhotoController struct {
	db *mongo.Database
}

func InitPhotoRoutes(db *mongo.Database, mux *http.ServeMux) {
	photoController := &PhotoController{db: db}

	mux.Handle("/photos", withAuth(http.HandlerFunc(photoController.GetPhotos)))
}

func (ctlr PhotoController) GetPhotos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

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

	collection := ctlr.db.Collection("photos")
	cursor, err := collection.Find(context.Background(), map[string]any{
		"timestamp": map[string]any{
			"$gte": time.Unix(startInt, 0),
			"$lte": time.Unix(endInt, 0),
		},
	})
	if err != nil {
		http.Error(w, "Failed to fetch photos", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var photos []Photo
	for cursor.Next(context.Background()) {
		var photo Photo
		if err := cursor.Decode(&photo); err != nil {
			http.Error(w, "Failed to decode photo", http.StatusInternalServerError)
			return
		}
		presignedURL, err := utils.GetPresignedURL(context.Background(), fmt.Sprintf("photos/%d.%s", photo.Timestamp.Unix(), photo.ImageType))
		if err != nil {
			http.Error(w, "Failed to get presigned URL", http.StatusInternalServerError)
			return
		}
		photo.PresignedURL = presignedURL
		photos = append(photos, photo)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(photos)
}
