package data

import (
	"context"
	"errors"
	"strings"

	"task_manager/models"

	"go.mongodb.org/mongo-driver/bson"
)

// User Operations

func (r *Repo) CreateUser(user *models.User) error {
	if r.Exists(models.AuthUser{Email: user.Email}) {
		return errors.New("email already exists")
	}
	if r.isMemory {
		r.users = append(r.users, *user)
		return nil
	}
	_, err := r.collection("user").InsertOne(context.Background(), user)
	return err
}

func (r *Repo) UpdateUser(id string, user models.AuthUser) error {
	if r.isMemory {
		for idx, u := range r.users {
			if u.ID == id {
				if user.Email != "" {
					r.users[idx].Email = user.Email
				}
				if user.Password != "" {
					r.users[idx].Password = user.Password
				}
				return nil
			}
		}
		return errors.New("user not found")
	}
	filter := bson.M{"id": id}
	updateData := bson.M{}
	if user.Email != "" {
		updateData["email"] = user.Email
	}
	if user.Password != "" {
		updateData["password"] = user.Password
	}
	update := bson.M{"$set": updateData}
	_, err := r.collection("user").UpdateOne(context.Background(), filter, update)
	return err
}

func (r *Repo) DeleteUser(id string) error {
	if r.isMemory {
		for idx, u := range r.users {
			if u.ID == id {
				r.users = append(r.users[:idx], r.users[idx+1:]...)
				r.tasks = filterTasks(r.tasks, func(t models.Task) bool {
					return t.Owner != id
				})
				return nil
			}
		}
		return errors.New("user not found")
	}
	_, err := r.collection("user").DeleteOne(context.Background(), bson.M{"id": id})
	if err != nil {
		return err
	}
	_, err = r.collection("tasks").DeleteMany(context.Background(), bson.M{"owner": id})
	return err
}

func (r *Repo) Exists(user models.AuthUser) bool {
	if r.isMemory {
		for _, u := range r.users {
			if strings.EqualFold(u.Email, user.Email) {
				return true
			}
		}
		return false
	}
	err := r.collection("user").FindOne(context.Background(), bson.M{
		"email": bson.M{"$regex": "^" + user.Email + "$", "$options": "i"},
	}).Err()
	return err == nil
}

func (r *Repo) GetUser(user models.AuthUser) (models.User, error) {
	if r.isMemory {
		for _, u := range r.users {
			if strings.EqualFold(u.Email, user.Email) {
				return u, nil
			}
		}
		return models.User{}, errors.New("no user found")
	}
	var userFound models.User
	err := r.collection("user").FindOne(context.Background(), bson.M{
		"email": bson.M{"$regex": "^" + user.Email + "$", "$options": "i"},
	}).Decode(&userFound)
	if err != nil {
		return models.User{}, errors.New("no user found")
	}
	return userFound, nil
}

func filterTasks(tasks []models.Task, predicate func(models.Task) bool) []models.Task {
	var filtered []models.Task
	for _, t := range tasks {
		if predicate(t) {
			filtered = append(filtered, t)
		}
	}
	return filtered
}