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

	"mqtt-streaming-server/domain"
	"mqtt-streaming-server/repository"
	"mqtt-streaming-server/utils"
)

type BrokerHandler struct {
	photoRepository  domain.PhotoRepository
	deviceRepository domain.DeviceRepository
	ocrClient        *gosseract.Client
}

func NewBrokerHandler(db *mongo.Database, ocrClient *gosseract.Client) BrokerHandler {
	return BrokerHandler{
		photoRepository:  repository.NewPhotoRepository(db),
		deviceRepository: repository.NewDeviceRepository(db),
		ocrClient:        ocrClient,
	}
}

func (b BrokerHandler) HandlePhoto(_ mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	// topci is photos/device_id
	deviceID := topic[len("photos/"):]
	ctx := context.Background()
	fmt.Println("Received message on topic:", msg.Topic())
	// get registered device
	device, err := b.deviceRepository.GetByID(ctx, deviceID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Printf("Device ID not found: %s\n", deviceID)
		} else {
			fmt.Printf("Failed to check device ID: %v\n", err)
		}
		return
	}
	fmt.Printf("Received photo from device: %s\n", device.DeviceName)
	body := msg.Payload()
	_, imageType, err := image.DecodeConfig(bytes.NewReader(body))
	if err != nil {
		fmt.Printf("Failed to decode image: %v\n", err)
		return
	}
	fmt.Printf("Image type: %s\n", imageType)

	// Extract text from image
	text, err := b.extractTextFromImage(body)
	if err != nil {
		fmt.Printf("Failed to extract text from image: %v\n", err)
		text = "OCR failed"
	}
	// UTC timestamp
	timestamp := time.Now().UTC()
	err = b.photoRepository.Save(ctx, &domain.Photo{
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
	// Check if device ID already exists
	_, err := b.deviceRepository.GetByID(ctx, deviceID)
	if err != nil && err != mongo.ErrNoDocuments {
		fmt.Printf("Failed to check device ID: %v\n", err)
		return
	}
	if err == mongo.ErrNoDocuments {
		// Device ID does not exist, insert it
		err = b.deviceRepository.Save(ctx, &domain.Device{
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
	err = b.deviceRepository.Update(ctx, deviceID, &domain.Device{
		DeviceID:     deviceID,
		DeviceName:   string(body),
		DeviceStatus: "active",
	})
	if err != nil {
		fmt.Printf("Failed to update device ID: %v\n", err)
		return
	}
	fmt.Printf("Device updated: %s\n", deviceID)
}

func (b BrokerHandler) DisconnectDevice(_ mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	// topic is disconnect/device_id
	deviceID := topic[len("device/id/"):]
	ctx := context.Background()
	fmt.Println("Received message on topic:", msg.Topic())
	message := string(msg.Payload())
	fmt.Printf("Received device disconnection: %s\n", message)
	// Check if message is a disconnect request
	if message != "Device disconnected" {
		fmt.Printf("Invalid disconnection message: %s\n", message)
		return
	}
	// Check if device ID exists
	device, err := b.deviceRepository.GetByID(ctx, deviceID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Printf("Device ID not found: %s\n", deviceID)
		} else {
			fmt.Printf("Failed to check device ID: %v\n", err)
		}
		return
	}
	// Update device status to inactive
	err = b.deviceRepository.Update(ctx, deviceID, &domain.Device{
		DeviceID:     deviceID,
		DeviceStatus: "inactive",
		DeviceName:   device.DeviceName,
	})
	if err != nil {
		fmt.Printf("Failed to update device ID: %v\n", err)
		return
	}
	fmt.Printf("Device disconnected: %s\n", deviceID)
}

func (b BrokerHandler) extractTextFromImage(imageData []byte) (string, error) {
	// Use the OCR client to extract text from the image
	b.ocrClient.SetImageFromBytes(imageData)
	text, err := b.ocrClient.Text()
	if err != nil {
		return "", fmt.Errorf("failed to extract text from image: %v", err)
	}
	return text, nil
}
