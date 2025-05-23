package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Photo struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Timestamp    time.Time          `json:"timestamp" bson:"timestamp"`
	ImageType    string             `json:"image_type" bson:"image_type"`
	PresignedURL string             `json:"presigned_url" bson:",omitempty"`
	DeviceID     string             `json:"device_id" bson:"device_id"`
	Text         string             `json:"text" bson:"text"`
}

type PhotoRepository interface {
	GetPhotos(ctx context.Context, filters map[string]any) ([]*Photo, error)
	Save(ctx context.Context, photo *Photo) error
}
