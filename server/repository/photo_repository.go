package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"mqtt-streaming-server/domain"
)

type photoRepository struct {
	db *mongo.Database
}

func NewPhotoRepository(db *mongo.Database) *photoRepository {
	return &photoRepository{db: db}
}

func (repo *photoRepository) GetPhotos(ctx context.Context, filters map[string]any) ([]*domain.Photo, error) {
	collection := repo.db.Collection("photos")
	photos := make([]*domain.Photo, 0)
	cursor, err := collection.Find(ctx, filters, &options.FindOptions{
		Sort: map[string]int{"timestamp": -1}, // Sort by timestamp in descending order
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var photo domain.Photo
		if err := cursor.Decode(&photo); err != nil {
			return nil, err
		}
		photos = append(photos, &photo)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return photos, nil
}

func (repo *photoRepository) Save(ctx context.Context, photo *domain.Photo) error {
	collection := repo.db.Collection("photos")
	_, err := collection.InsertOne(ctx, photo)
	return err
}
