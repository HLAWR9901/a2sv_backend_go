package repositories

import (
	"context"
	"task_manager/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const taskCollection = "tasks"

type TaskRepository struct {
	Collection *mongo.Collection
}

// Init Task Repo
func NewTaskRepository(repo *Repository) domain.ITaskRepository {
	return &TaskRepository{Collection: repo.Database.Collection(taskCollection)}
}

// Task DB Operations
func (tr *TaskRepository) Create(ctx context.Context, task *domain.Task) error {
	_, err := tr.Collection.InsertOne(ctx, task)
	return err
}

func (tr *TaskRepository) Update(ctx context.Context, task *domain.Task) error {
	filter := bson.M{"_id": task.ID}
	update := bson.M{"$set": bson.M{
		"title":       task.Title,
		"description": task.Description,
		"status":      task.Status,
		"updated_at":  time.Now(),
	}}
	_, err := tr.Collection.UpdateOne(ctx, filter, update)
	return err
}

func (tr *TaskRepository) Delete(ctx context.Context, id string) error {
	_, err := tr.Collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (tr *TaskRepository) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	var task domain.Task
	err := tr.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (tr *TaskRepository) GetByUser(ctx context.Context, userID string) ([]domain.Task, error) {
	cursor, err := tr.Collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	// Iterate over the cursor
	var tasks []domain.Task
	for cursor.Next(ctx) {
		var task domain.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}