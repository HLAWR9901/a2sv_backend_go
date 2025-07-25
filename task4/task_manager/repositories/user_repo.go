package repositories

import (
	"context"
	"task_manager/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const userCollection = "users"

type UserRepository struct {
	Collection *mongo.Collection
}

// Init UserRepo
func NewUserRepository(repo *Repository) domain.IUserRepository {
	return &UserRepository{Collection: repo.Database.Collection(userCollection)}
}

// User DB Operations
func (ur *UserRepository) Create(ctx context.Context, user *domain.User) error {
	_, err := ur.Collection.InsertOne(ctx, user)
	return err
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := ur.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := ur.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := ur.Collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}