package broker

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/mongo"

	"mqtt-streaming-server/routes"
	"mqtt-streaming-server/utils"
)

type BrokerHandler struct {
	db *mongo.Database
}

func NewBrokerHandler(db *mongo.Database) BrokerHandler {
	return BrokerHandler{
		db: db,
	}
}

func (b BrokerHandler) HandlePhoto(_ mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	// topci is photos/device_id
	deviceID := topic[len("photos/"):]
	ctx := context.Background()
	fmt.Println("Received message on topic:", msg.Topic())
	// get registered device
	collection := b.db.Collection("devices")
	var result routes.Device
	err := collection.FindOne(ctx, map[string]string{"device_id": deviceID}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Printf("Device ID not found: %s\n", deviceID)
		} else {
			fmt.Printf("Failed to check device ID: %v\n", err)
		}
		return
	}
	fmt.Printf("Received photo from device: %s\n", result.DeviceName)
	body := msg.Payload()
	_, imageType, err := image.DecodeConfig(bytes.NewReader(body))
	if err != nil {
		fmt.Printf("Failed to decode image: %v\n", err)
		return
	}
	fmt.Printf("Image type: %s\n", imageType)
	// UTC timestamp
	timestamp := time.Now().UTC()
	collection = b.db.Collection("photos")
	_, err = collection.InsertOne(ctx, routes.Photo{
		ImageType: imageType,
		Timestamp: timestamp,
		DeviceID:  deviceID,
	})
	if err != nil {
		fmt.Printf("Failed to insert photo into MongoDB: %v\n", err)
		return
	}
	// upload to S3
	keyName := fmt.Sprintf("photos/%d.%s", timestamp.Unix(), imageType)
	err = utils.UploadToS3(ctx, body, imageType, keyName)
	if err != nil {
		fmt.Printf("Failed to upload photo to S3: %v\n", err)
		return
	}
	fmt.Printf("Photo uploaded to S3 with key: %s\n", keyName)
}

func (b BrokerHandler) RegisterDevice(_ mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	// topic is register/device_id
	deviceID := topic[len("register/"):]
	ctx := context.Background()
	fmt.Println("Received message on topic:", msg.Topic())
	body := msg.Payload()
	fmt.Printf("Received device registration: %s\n", body)
	collection := b.db.Collection("devices")
	// Check if device ID already exists
	var result routes.Device
	err := collection.FindOne(ctx, map[string]string{"device_id": deviceID}).Decode(&result)
	if err != nil && err != mongo.ErrNoDocuments {
		fmt.Printf("Failed to check device ID: %v\n", err)
		return
	}
	if err == mongo.ErrNoDocuments {
		// Device ID does not exist, insert it
		_, err = collection.InsertOne(ctx, routes.Device{
			DeviceID:     deviceID,
			DeviceName:   string(body),
			DeviceStatus: "active",
		})
		if err != nil {
			fmt.Printf("Failed to insert device ID: %v\n", err)
			return
		}
		fmt.Printf("Device registered: %s\n", deviceID)
		return
	}
	// Device ID already exists, update it
	_, err = collection.UpdateOne(ctx, map[string]string{"device_id": deviceID}, map[string]any{
		"$set": map[string]any{
			"device_name":   string(body),
			"device_status": "active",
		},
	})
	if err != nil {
		fmt.Printf("Failed to update device ID: %v\n", err)
		return
	}
	fmt.Printf("Device updated: %s\n", deviceID)
}
