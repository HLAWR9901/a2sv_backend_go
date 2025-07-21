package data

import (
    "context"
    "errors"
    "task_manager/models"

    "github.com/google/uuid"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

type Repo struct {
    Client   *mongo.Client
    isMemory bool
    tasks    []models.Task
}

func NewRepo(client *mongo.Client, isMemory bool) *Repo {
    return &Repo{Client: client, isMemory: isMemory, tasks: []models.Task{}}
}

func (r *Repo) collection(coll string) *mongo.Collection {
    return r.Client.Database("task_manager_db").Collection(coll)
}

func (r *Repo) Create(task *models.Task) error {
    if r.isMemory {
        r.tasks = append(r.tasks, *task)
        return nil
    }
    _, err := r.collection("tasks").InsertOne(context.Background(), task)
    return err
}

func (r *Repo) Update(id string, task models.Task) error {
    u, err := uuid.Parse(id)
    if err != nil {
        return errors.New("invalid UUID")
    }
    if r.isMemory {
        for i, t := range r.tasks {
            if t.ID == u {
                if task.Name != "" {
                    r.tasks[i].Name = task.Name
                }
                if task.Description != nil {
                    r.tasks[i].Description = task.Description
                }
                if task.Status != "" {
                    r.tasks[i].Status = task.Status
                }
                if task.Priority != "" {
                    r.tasks[i].Priority = task.Priority
                }
                if task.DueDate != nil {
                    r.tasks[i].DueDate = task.DueDate
                }
                r.tasks[i].UpdatedAt = task.UpdatedAt
                return nil
            }
        }
        return errors.New("task not found")
    }

    filter := bson.M{"id": u}
    updateData := bson.M{"updated_at": task.UpdatedAt}
    if task.Name != "" {
        updateData["name"] = task.Name
    }
    if task.Description != nil {
        updateData["description"] = task.Description
    }
    if task.Status != "" {
        updateData["status"] = task.Status
    }
    if task.Priority != "" {
        updateData["priority"] = task.Priority
    }
    if task.DueDate != nil {
        updateData["due_date"] = task.DueDate
    }

    res, err := r.collection("tasks").UpdateOne(context.Background(), filter, bson.M{"$set": updateData})
    if err != nil {
        return err
    }
    if res.MatchedCount == 0 {
        return errors.New("task not found")
    }
    return nil
}

func (r *Repo) Delete(id string) error {
    u, err := uuid.Parse(id)
    if err != nil {
        return errors.New("invalid UUID")
    }
    if r.isMemory {
        for i, t := range r.tasks {
            if t.ID == u {
                r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)
                return nil
            }
        }
        return errors.New("task not found")
    }
    res, err := r.collection("tasks").DeleteOne(context.Background(), bson.M{"id": u})
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
    if err := cursor.Err(); err != nil {
        return nil, err
    }
    return tasks, nil
}

func (r *Repo) GetById(id string) (*models.Task, error) {
    u, err := uuid.Parse(id)
    if err != nil {
        return nil, errors.New("invalid UUID")
    }
    if r.isMemory {
        for _, t := range r.tasks {
            if t.ID == u {
                return &t, nil
            }
        }
        return nil, errors.New("task not found")
    }
    var t models.Task
    err = r.collection("tasks").FindOne(context.Background(), bson.M{"id": u}).Decode(&t)
    if err == mongo.ErrNoDocuments {
        return nil, errors.New("task not found")
    }
    if err != nil {
        return nil, err
    }
    return &t, nil
}