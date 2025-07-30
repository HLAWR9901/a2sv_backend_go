package repositories

import "go.mongodb.org/mongo-driver/mongo"

// Repository
type Repository struct {
	Database *mongo.Database
}

func NewRepository(client *mongo.Client, dbName string) *Repository {
	return &Repository{Database: client.Database(dbName)}
}