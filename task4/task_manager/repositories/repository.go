package repositories

import "go.mongodb.org/mongo-driver/mongo"

// Repository holds the database connection
type Repository struct {
	Database *mongo.Database
}

func NewRepository(client *mongo.Client, dbName string) *Repository {
	return &Repository{Database: client.Database(dbName)}
}