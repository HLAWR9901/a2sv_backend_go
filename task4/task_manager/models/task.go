package models

import (
	"time"
	"github.com/google/uuid"
)
type State string

const (
	Pending State = "pending"
	InProgress State = "inprogress"
	Completed State = "completed"
)

type Importance string

const (
	High Importance = "high"
	Medium Importance = "medium"
	Low Importance = "low"
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
