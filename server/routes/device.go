package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/mongo"
)

type Device struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	DeviceID     string `json:"device_id" bson:"device_id"`
	DeviceName   string `json:"device_name" bson:"device_name"`
	DeviceStatus string `json:"device_status" bson:"device_status"`
}

type DeviceController struct {
	db         *mongo.Database
	mqttClient mqtt.Client
}

func InitDeviceRoutes(db *mongo.Database, mqttClient mqtt.Client, mux *http.ServeMux) {
	deviceController := &DeviceController{db: db, mqttClient: mqttClient}

	mux.Handle("/devices", withAuth(http.HandlerFunc(deviceController.GetDevices)))
	mux.Handle("/devices/switch", withAuth(http.HandlerFunc(deviceController.SwitchDeviceMode)))
}

func (ctlr DeviceController) SwitchDeviceMode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	// Check if the user is authorized
	if ctx.Value("role") != "admin" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var device struct {
		ID   string `json:"id"`
		Mode string `json:"mode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	topic := fmt.Sprintf("setup/%s", device.ID)
	if token := ctlr.mqttClient.Publish(topic, 0, false, "set "+device.Mode); token.Wait() && token.Error() != nil {
		http.Error(w, "Failed to publish message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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
