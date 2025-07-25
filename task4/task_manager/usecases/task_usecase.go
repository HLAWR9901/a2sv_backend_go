package usecases

import (
	"context"
	"errors"
	"task_manager/domain"
	"task_manager/infrastructure/config"
	"time"

	"github.com/google/uuid"
)

type TaskUsecase struct {
	repo domain.ITaskRepository
}

// Init Task Usecase
func NewTaskUsecase(repo domain.ITaskRepository) *TaskUsecase {
	return &TaskUsecase{repo: repo}
}

// Tast Usecase application logic
func (tu *TaskUsecase) Create(ctx context.Context, task *domain.Task) error {
	ctx, cancel := context.WithTimeout(ctx, config.TaskCreationTimeout)
	defer cancel()

	if task.Title == "" || task.Status == "" {
		return errors.New(domain.ErrInvalidPayload)
	}

	task.ID = uuid.NewString()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	return tu.repo.Create(ctx, task)
}

func (tu *TaskUsecase) UpdateTask(ctx context.Context, userID, taskID string, task domain.Task) (*domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, config.TaskUpdateTimeout)
	defer cancel()

	existing, err := tu.repo.GetByID(ctx, taskID)
	if err != nil {
		return nil, errors.New(domain.ErrTaskNotFound)
	}

	if existing.UserID != userID {
		return nil, errors.New(domain.ErrAccessDenied)
	}

	// Update fields
	if task.Title != "" {
		existing.Title = task.Title
	}
	if task.Description != "" {
		existing.Description = task.Description
	}
	if task.Status != "" {
		existing.Status = task.Status
	}
	existing.UpdatedAt = time.Now()

	if err := tu.repo.Update(ctx, existing); err != nil {
		return nil, errors.New(domain.ErrInternalServer)
	}

	return existing, nil
}

func (tu *TaskUsecase) DeleteTask(ctx context.Context, userID, taskID string) error {
	ctx, cancel := context.WithTimeout(ctx, config.TaskDeleteTimeout)
	defer cancel()

	task, err := tu.repo.GetByID(ctx, taskID)
	if err != nil {
		return errors.New(domain.ErrTaskNotFound)
	}

	if task.UserID != userID {
		return errors.New(domain.ErrAccessDenied)
	}

	return tu.repo.Delete(ctx, taskID)
}

func (tu *TaskUsecase) GetTaskByID(ctx context.Context, userID, taskID string) (*domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, config.GetTaskByIdTimeout)
	defer cancel()

	task, err := tu.repo.GetByID(ctx, taskID)
	if err != nil {
		return nil, errors.New(domain.ErrTaskNotFound)
	}

	if task.UserID != userID {
		return nil, errors.New(domain.ErrAccessDenied)
	}

	return task, nil
}

func (tu *TaskUsecase) GetTasksByUser(ctx context.Context, userID string) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, config.GetTaskByUserTimeout)
	defer cancel()

	tasks, err := tu.repo.GetByUser(ctx, userID)
	if err != nil {
		return nil, errors.New(domain.ErrInternalServer)
	}

	if len(tasks) == 0 {
		return nil, errors.New(domain.ErrNoTaskByUser)
	}

	return tasks, nil
}