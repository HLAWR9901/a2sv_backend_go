package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"task_manager/domain"
	"task_manager/mocks"
)

func TestTaskUsecase_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	usecase := NewTaskUsecase(mockRepo)
	ctx := context.Background()

	task := &domain.Task{Title: "Test Task", Status: domain.Pending, UserID: "1"}

	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	err := usecase.Create(ctx, task)

	assert.NoError(t, err)
	assert.NotEmpty(t, task.ID)
	assert.NotZero(t, task.CreatedAt)
	assert.NotZero(t, task.UpdatedAt)

	mockRepo.AssertExpectations(t)
}

func TestTaskUsecase_Create_InvalidPayload(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	usecase := NewTaskUsecase(mockRepo)
	ctx := context.Background()

	task1 := &domain.Task{Title: "", Status: domain.Pending, UserID: "1"}
	err1 := usecase.Create(ctx, task1)
	assert.Error(t, err1)
	assert.Equal(t, domain.ErrInvalidPayload, err1.Error())

	task2 := &domain.Task{Title: "Test Task", Status: "", UserID: "1"}
	err2 := usecase.Create(ctx, task2)
	assert.Error(t, err2)
	assert.Equal(t, domain.ErrInvalidPayload, err2.Error())

	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestTaskUsecase_Create_RepoError(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	usecase := NewTaskUsecase(mockRepo)
	ctx := context.Background()

	task := &domain.Task{Title: "Test Task", Status: domain.Pending, UserID: "1"}
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))
	err := usecase.Create(ctx, task)

	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())

	assert.NotEmpty(t, task.ID)
	assert.NotZero(t, task.CreatedAt)
	assert.NotZero(t, task.UpdatedAt)

	mockRepo.AssertExpectations(t)
}

func TestTaskUsecase_UpdateTask_Success(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	usecase := NewTaskUsecase(mockRepo)
	ctx := context.Background()

	userID := "1"
	taskID := "task1"

	existingTask := &domain.Task{ID: taskID, UserID: userID, Title: "Old Title", Status: domain.Pending}
	updatedTask := domain.Task{Title: "New Title", Description: "New Desc", Status: domain.Completed}

	mockRepo.On("GetByID", mock.Anything, taskID).Return(existingTask, nil)
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	result, err := usecase.UpdateTask(ctx, userID, taskID, updatedTask)

	assert.NoError(t, err)
	assert.Equal(t, existingTask, result)
	assert.Equal(t, "New Title", result.Title)
	assert.Equal(t, "New Desc", result.Description)
	assert.Equal(t, domain.Completed, result.Status)

	mockRepo.AssertExpectations(t)
}

func TestTaskUsecase_UpdateTask_TaskNotFound(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	usecase := NewTaskUsecase(mockRepo)
	ctx := context.Background()

	userID := "1"
	taskID := "task1"
	updatedTask := domain.Task{Title: "New Title"}
	mockRepo.On("GetByID", mock.Anything, taskID).Return((*domain.Task)(nil), errors.New("not found"))
	result, err := usecase.UpdateTask(ctx, userID, taskID, updatedTask)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrTaskNotFound, err.Error())
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestTaskUsecase_UpdateTask_AccessDenied(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	usecase := NewTaskUsecase(mockRepo)
	ctx := context.Background()

	userID := "1"
	taskID := "task1"
	updatedTask := domain.Task{Title: "New Title"}

	existingTask := &domain.Task{ID: taskID, UserID: "2", Title: "Old Title", Status: domain.Pending}
	mockRepo.On("GetByID", mock.Anything, taskID).Return(existingTask, nil)
	result, err := usecase.UpdateTask(ctx, userID, taskID, updatedTask)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrAccessDenied, err.Error())
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestTaskUsecase_DeleteTask_Success(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	usecase := NewTaskUsecase(mockRepo)
	ctx := context.Background()

	userID := "1"
	taskID := "task1"

	existingTask := &domain.Task{ID: taskID, UserID: userID}
	mockRepo.On("GetByID", mock.Anything, taskID).Return(existingTask, nil)
	mockRepo.On("Delete", mock.Anything, taskID).Return(nil)

	err := usecase.DeleteTask(ctx, userID, taskID)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestTaskUsecase_DeleteTask_TaskNotFound(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	usecase := NewTaskUsecase(mockRepo)
	ctx := context.Background()

	userID := "1"
	taskID := "task1"
	mockRepo.On("GetByID", mock.Anything, taskID).Return((*domain.Task)(nil), errors.New("not found"))
	err := usecase.DeleteTask(ctx, userID, taskID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrTaskNotFound, err.Error())

	mockRepo.AssertExpectations(t)
}

func TestTaskUsecase_DeleteTask_AccessDenied(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	usecase := NewTaskUsecase(mockRepo)
	ctx := context.Background()

	userID := "1"
	taskID := "task1"

	existingTask := &domain.Task{ID: taskID, UserID: "2"}
	mockRepo.On("GetByID", mock.Anything, taskID).Return(existingTask, nil)
	err := usecase.DeleteTask(ctx, userID, taskID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrAccessDenied, err.Error())

	mockRepo.AssertExpectations(t)
}

func TestTaskUsecase_GetTaskByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	usecase := NewTaskUsecase(mockRepo)
	ctx := context.Background()

	userID := "1"
	taskID := "task1"
	existingTask := &domain.Task{ID: taskID, UserID: userID, Title: "Test Task", Status: domain.Pending}
	mockRepo.On("GetByID", mock.Anything, taskID).Return(existingTask, nil)
	result, err := usecase.GetTaskByID(ctx, userID, taskID)

	assert.NoError(t, err)
	assert.Equal(t, existingTask, result)

	mockRepo.AssertExpectations(t)
}

func TestTaskUsecase_GetTaskByID_TaskNotFound(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	usecase := NewTaskUsecase(mockRepo)
	ctx := context.Background()

	userID := "1"
	taskID := "task1"
	mockRepo.On("GetByID", mock.Anything, taskID).Return((*domain.Task)(nil), errors.New("not found"))
	result, err := usecase.GetTaskByID(ctx, userID, taskID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrTaskNotFound, err.Error())
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestTaskUsecase_GetTaskByID_AccessDenied(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	usecase := NewTaskUsecase(mockRepo)
	ctx := context.Background()

	userID := "1"
	taskID := "task1"
	existingTask := &domain.Task{ID: taskID, UserID: "2"}
	mockRepo.On("GetByID", mock.Anything, taskID).Return(existingTask, nil)
	result, err := usecase.GetTaskByID(ctx, userID, taskID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrAccessDenied, err.Error())
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestTaskUsecase_GetTasksByUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	usecase := NewTaskUsecase(mockRepo)
	ctx := context.Background()

	userID := "1"
	tasks := []domain.Task{
		{ID: "task1", UserID: userID, Title: "Task 1", Status: domain.Pending},
		{ID: "task2", UserID: userID, Title: "Task 2", Status: domain.Completed},
	}
	mockRepo.On("GetByUser", mock.Anything, userID).Return(tasks, nil)
	result, err := usecase.GetTasksByUser(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, tasks, result)

	mockRepo.AssertExpectations(t)
}

func TestTaskUsecase_GetTasksByUser_NoTasks(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	usecase := NewTaskUsecase(mockRepo)
	ctx := context.Background()

	userID := "1"
	mockRepo.On("GetByUser", mock.Anything, userID).Return([]domain.Task{}, nil)
	result, err := usecase.GetTasksByUser(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrNoTaskByUser, err.Error())
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestTaskUsecase_GetTasksByUser_RepoError(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	usecase := NewTaskUsecase(mockRepo)
	ctx := context.Background()

	userID := "1"
	mockRepo.On("GetByUser", mock.Anything, userID).Return([]domain.Task{}, errors.New("db error"))
	result, err := usecase.GetTasksByUser(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInternalServer, err.Error())
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}