package data

import (
	"context"
	"errors"

	"task_manager/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repo General
type Repo struct {
	Client   *mongo.Client
	isMemory bool
	tasks    []models.Task
	users    []models.User
}

func NewRepo(client *mongo.Client, isMemory bool) *Repo {
	return &Repo{Client: client, isMemory: isMemory, tasks: []models.Task{}, users: []models.User{}}
}

func (r *Repo) collection(coll string) *mongo.Collection {
	return r.Client.Database("task_manager_db").Collection(coll)
}

//Task Operations

func (r *Repo) Create(task *models.Task) error {
	if r.isMemory {
		r.tasks = append(r.tasks, *task)
		return nil
	}
	_, err := r.collection("tasks").InsertOne(context.Background(), task)
	return err
}

func (r *Repo) Update(id string, task models.Task) error {
	if r.isMemory {
		for i, t := range r.tasks {
			if t.ID == id {
				r.tasks[i] = task
				return nil
			}
		}
		return errors.New("task not found")
	}
	filter := bson.M{"id": id}
	update := bson.M{"$set": task}
	res, err := r.collection("tasks").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("task not found")
	}
	return nil
}

func (r *Repo) Delete(id string) error {
	if r.isMemory {
		for i, t := range r.tasks {
			if t.ID == id {
				r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)
				return nil
			}
		}
		return errors.New("task not found")
	}
	res, err := r.collection("tasks").DeleteOne(context.Background(), bson.M{"id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("task not found")
	}
	return nil
}

func (r *Repo) GetAll() ([]models.Task, error) {
	if r.isMemory {
		return r.tasks, nil
	}
	cursor, err := r.collection("tasks").Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var tasks []models.Task
	for cursor.Next(context.Background()) {
		var t models.Task
		if err := cursor.Decode(&t); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *Repo) GetById(id string) (*models.Task, error) {
	if r.isMemory {
		for _, t := range r.tasks {
			if t.ID == id {
				return &t, nil
			}
		}
		return nil, errors.New("task not found")
	}
	var t models.Task
	err := r.collection("tasks").FindOne(context.Background(), bson.M{"id": id}).Decode(&t)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("task not found")
	}
	return &t, err
}

func (r *Repo) GetByOwner(owner string) ([]models.Task, error) {
	if r.isMemory {
		var tasks []models.Task
		for _, t := range r.tasks {
			if t.Owner == owner {
				tasks = append(tasks, t)
			}
		}
		return tasks, nil
	}
	cursor, err := r.collection("tasks").Find(context.Background(), bson.M{"owner": owner})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var tasks []models.Task
	for cursor.Next(context.Background()) {
		var t models.Task
		if err := cursor.Decode(&t); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *Repo) CleanAll() error {
	if r.isMemory {
		r.tasks = []models.Task{}
		return nil
	}
	_,err:=r.collection("tasks").DeleteMany(context.TODO(), bson.M{})
	return err
}