package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/mongo"

	"mqtt-streaming-server/domain"
	"mqtt-streaming-server/repository"
)

type DeviceController struct {
	DeviceRepository domain.DeviceRepository
	mqttClient       mqtt.Client
}

func InitDeviceRoutes(db *mongo.Database, mqttClient mqtt.Client, mux *http.ServeMux) {
	deviceController := &DeviceController{
		DeviceRepository: repository.NewDeviceRepository(db),
		mqttClient:       mqttClient,
	}

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
	devices, err := ctlr.DeviceRepository.GetAllDevices(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch devices", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}
