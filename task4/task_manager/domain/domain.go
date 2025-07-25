package domain

import (
	"time"
	"context"
)

// User model
type User struct {
	ID       string `bson:"_id"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

// Task model
type State string

const (
	Pending     State = "pending"
	InProgress  State = "inprogress"
	Completed   State = "completed"
)

type Task struct {
	ID          string    `bson:"_id"`
	Title       string    `bson:"title"`
	Description string    `bson:"description"`
	Status      State     `bson:"status"`
	UserID      string    `bson:"user_id"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

// Repository interfaces
type IUserRepository interface {
	Create(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	Delete(ctx context.Context, id string) error
}

type ITaskRepository interface {
	Create(ctx context.Context, task *Task) error
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Task, error)
	GetByUser(ctx context.Context, userID string) ([]Task, error)
}

// Service interfaces
type IPasswordService interface {
	Hash(password string) (string, error)
	Verify(hash, password string) (bool, error)
}

type IAuthService interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(token string) (string, error)
}

// Usecase interfaces
type IUserUsecase interface {
	Create(ctx context.Context, user *User) (*User, error)
	Login(ctx context.Context, user *User) (*User, string, error)
	Logout(ctx context.Context) (string, error)
	DeleteUserAndTasks(ctx context.Context, userID string) error
	VerifyPassword(ctx context.Context, userID, password string) error
}

type ITaskUsecase interface {
	Create(ctx context.Context, task *Task) error
	UpdateTask(ctx context.Context, userID, taskID string, task Task) (*Task, error)
	DeleteTask(ctx context.Context, userID, taskID string) error
	GetTaskByID(ctx context.Context, userID, taskID string) (*Task, error)
	GetTasksByUser(ctx context.Context, userID string) ([]Task, error)
}

// Error constants
const (
	ErrUserExists      = "user already exists"
	ErrInternalServer  = "internal server error"
	ErrInvalidPayload  = "invalid request payload"
	ErrUserNotFound    = "user not found"
	ErrHashingPassword = "failed to hash password"
	ErrTaskNotFound    = "task not found"
	ErrNoTaskByUser    = "user has no tasks"
	ErrAccessDenied    = "access denied"
	ErrInvalidToken    = "invalid token"
)