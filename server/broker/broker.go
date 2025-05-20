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
	ctx := context.Background()
	fmt.Println("Received message on topic:", msg.Topic())
	if msg.Topic() != "photos" {
		fmt.Printf("Invalid topic: %s", msg.Topic())
		return
	}
	body := msg.Payload()
	_, imageType, err := image.DecodeConfig(bytes.NewReader(body))
	if err != nil {
		fmt.Printf("Failed to decode image: %v\n", err)
		return
	}
	fmt.Printf("Image type: %s\n", imageType)
	// UTC timestamp
	timestamp := time.Now().UTC()
	collection := b.db.Collection("photos")
	_, err = collection.InsertOne(ctx, routes.Photo{
		ImageType: imageType,
		Timestamp: timestamp,
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
