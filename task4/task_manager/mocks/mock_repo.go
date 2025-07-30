package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"task_manager/domain"
)

// User Mocks
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockPasswordService struct {
	mock.Mock
}

func (m *MockPasswordService) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordService) Verify(hashedPassword, password string) (bool, error) {
	args := m.Called(hashedPassword, password)
	return args.Bool(0), args.Error(1)
}

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) GenerateToken(userID string) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ValidateToken(tokenStr string) (string, error){
	args := m.Called(tokenStr)
	return args.String(0),args.Error(1)
}

// Task Mocks
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, task *domain.Task) error{
	args := m.Called(ctx,task)
	return args.Error(0)
}

func (m *MockTaskRepository) Update(ctx context.Context, task *domain.Task) error{
	args := m.Called(ctx,task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(ctx context.Context, id string) error{
	args := m.Called(ctx,id)
	return args.Error(0)
}

func (m *MockTaskRepository) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockTaskRepository) GetByUser(ctx context.Context, userID string) ([]domain.Task, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.Task), args.Error(1)
}