package models

import (
	"time"
	"github.com/google/uuid"
)
type State string

const (
	Pending State = "Pending"
	InProgress State = "Inprogress"
	Completed State = "Completed"
)

type Importance string

const (
	High Importance = "High"
	Medium Importance = "Medium"
	Low Importance = "Low"
)

type BaseModel struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`

}
type Task struct {
	BaseModel
	Name string `json:"name"`
	Description string `json:"description"`
	Status State `json:"status"`
	Priority Importance `json:"priority"`
	DueDate time.Time `json:"duedate"`
}

func NewTask(
  name, desc string, status State,
  priority Importance,
  dueDate time.Time, createdAt time.Time,
) *Task {
  return &Task{
    BaseModel: BaseModel{
      ID: uuid.New(),
      CreatedAt: createdAt,
      UpdatedAt: createdAt,
    },
    Name: name,
    Description: desc,
    Status: status,
    Priority: priority,
    DueDate: dueDate,
  }
}
