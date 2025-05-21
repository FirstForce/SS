package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type Device struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	DeviceID     string `json:"device_id" bson:"device_id"`
	DeviceName   string `json:"device_name" bson:"device_name"`
	DeviceStatus string `json:"device_status" bson:"device_status"`
}

type DeviceController struct {
	db *mongo.Database
}

func InitDeviceRoutes(db *mongo.Database, mux *http.ServeMux) {
	deviceController := &DeviceController{db: db}

	mux.Handle("/devices", withAuth(http.HandlerFunc(deviceController.GetDevices)))
}

func (ctlr DeviceController) GetDevices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	// Check if the user is authorized
	if ctx.Value("role") != "admin" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch devices from the database
	collection := ctlr.db.Collection("devices")
	cursor, err := collection.Find(ctx, map[string]any{})
	if err != nil {
		fmt.Println("Failed to fetch devices:", err)
		http.Error(w, "Failed to fetch devices", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var devices []Device
	if err := cursor.All(ctx, &devices); err != nil {
		http.Error(w, "Failed to decode devices", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}
