package models

import (
	"errors"
	"time"
)

type State string
type Importance string

const (
	Pending    State = "Pending"
	InProgress State = "In Progress"
	Completed  State = "Completed"
)

const (
	Low    Importance = "Low"
	Medium Importance = "Medium"
	High   Importance = "High"
)

type Task struct {
	ID          string      `json:"id" bson:"id"`
	Name        string      `json:"name" bson:"name"`
	Description string      `json:"description" bson:"description"`
	Priority    Importance  `json:"priority" bson:"priority"`
	Status      State       `json:"status" bson:"status"`
	Owner       string      `json:"owner" bson:"owner"`
	DueDate     *time.Time  `json:"dueDate,omitempty" bson:"dueDate,omitempty"`
	CreatedAt   time.Time   `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt" bson:"updatedAt"`
}

func (t *Task) Validate() error {
	if t.Name == "" {
		return errors.New("name cannot be empty")
	}
	switch t.Priority {
	case Low, Medium, High:
	default:
		return errors.New("priority must be Low, Medium, or High")
	}
	switch t.Status {
	case Pending, InProgress, Completed:
	default:
		return errors.New("status must be Pending, In Progress, or Completed")
	}
	return nil
}