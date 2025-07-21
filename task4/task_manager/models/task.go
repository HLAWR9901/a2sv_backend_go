package models

import (
    "errors"
    "time"

    "github.com/google/uuid"
)

type State string
const (
    Pending    State = "pending"
    InProgress State = "inprogress"
    Completed  State = "completed"
)

type Importance string
const (
    High   Importance = "high"
    Medium Importance = "medium"
    Low    Importance = "low"
)

var (
    ValidStates = map[State]bool{
        Pending:    true,
        InProgress: true,
        Completed:  true,
    }
    ValidPriorities = map[Importance]bool{
        High:   true,
        Medium: true,
        Low:    true,
    }
)

type BaseModel struct {
    ID          uuid.UUID  `bson:"id" json:"id"`
    CreatedAt   time.Time  `bson:"created_at" json:"created_at"`
    UpdatedAt   time.Time  `bson:"updated_at" json:"updated_at"`
    CompletedAt *time.Time `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
}

type Task struct {
    BaseModel  `bson:",inline"`
    Name        string       `bson:"name" json:"name"`
    Description *string      `bson:"description,omitempty" json:"description,omitempty"`
    Status      State        `bson:"status" json:"status"`
    Priority    Importance   `bson:"priority" json:"priority"`
    DueDate     *time.Time   `bson:"due_date,omitempty" json:"due_date,omitempty"`
}

func NewTask(name, desc string, status State, priority Importance, dueDate *time.Time, createdAt time.Time) *Task {
    t := &Task{
        BaseModel: BaseModel{
            ID:        uuid.New(),
            CreatedAt: createdAt,
            UpdatedAt: createdAt,
        },
        Name:     name,
        Status:   status,
        Priority: priority,
        DueDate:  dueDate,
    }
    if desc != "" {
        t.Description = &desc
    }
    return t
}

func (t *Task) Validate() error {
    if t.Name == "" {
        return errors.New("name is required")
    }
    if !ValidStates[t.Status] {
        return errors.New("invalid status")
    }
    if !ValidPriorities[t.Priority] {
        return errors.New("invalid priority")
    }
    if t.DueDate != nil && !t.DueDate.IsZero() && t.DueDate.Before(t.CreatedAt) {
        return errors.New("due date is in the past")
    }
    return nil
}