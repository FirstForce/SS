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
	"github.com/otiai10/gosseract/v2"
	"go.mongodb.org/mongo-driver/mongo"

	"mqtt-streaming-server/routes"
	"mqtt-streaming-server/utils"
)

type BrokerHandler struct {
	db        *mongo.Database
	ocrClient *gosseract.Client
}

func NewBrokerHandler(db *mongo.Database, ocrClient *gosseract.Client) BrokerHandler {
	return BrokerHandler{
		db:        db,
		ocrClient: ocrClient,
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

	// Extract text from image
	text, err := b.ExtractTextFromImage(body)
	if err != nil {
		fmt.Printf("Failed to extract text from image: %v\n", err)
		text = "OCR failed"
	}
	// UTC timestamp
	timestamp := time.Now().UTC()
	collection = b.db.Collection("photos")
	_, err = collection.InsertOne(ctx, routes.Photo{
		ImageType: imageType,
		Timestamp: timestamp,
		DeviceID:  deviceID,
		Text:      text,
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

func (b BrokerHandler) ExtractTextFromImage(imageData []byte) (string, error) {
	// Use the OCR client to extract text from the image
	b.ocrClient.SetImageFromBytes(imageData)
	text, err := b.ocrClient.Text()
	if err != nil {
		return "", fmt.Errorf("failed to extract text from image: %v", err)
	}
	return text, nil
}
