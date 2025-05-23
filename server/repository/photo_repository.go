package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

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
	var photos []*domain.Photo
	cursor, err := collection.Find(ctx, filters)
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
