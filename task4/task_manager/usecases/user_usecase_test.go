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

func TestUserUsecase_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPwSvc := new(mocks.MockPasswordService)
	mockAuthSvc := new(mocks.MockAuthService)

	usecase := NewUserUsecase(mockRepo, mockPwSvc, mockAuthSvc)
	ctx := context.Background()
	user := &domain.User{Email: "test@example.com", Password: "password123"}
	hashedPassword := "hashedpassword"

	mockPwSvc.On("Hash", "password123").Return(hashedPassword, nil)
	mockRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return((*domain.User)(nil), nil)
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	result, err := usecase.Create(ctx, user)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, hashedPassword, result.Password)
	assert.NotEmpty(t, result.ID)

	mockPwSvc.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecase_Create_UserExists(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPwSvc := new(mocks.MockPasswordService)
	mockAuthSvc := new(mocks.MockAuthService)

	usecase := NewUserUsecase(mockRepo, mockPwSvc, mockAuthSvc)
	ctx := context.Background()
	user := &domain.User{Email: "test@example.com", Password: "password123"}
	existingUser := &domain.User{ID: "1", Email: "test@example.com", Password: "hashed"}

	mockRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(existingUser, nil)
	result, err := usecase.Create(ctx, user)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrUserExists, err.Error())
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
	mockPwSvc.AssertNotCalled(t, "Hash", mock.Anything)
}

func TestUserUsecase_Create_InvalidPayload(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPwSvc := new(mocks.MockPasswordService)
	mockAuthSvc := new(mocks.MockAuthService)

	usecase := NewUserUsecase(mockRepo, mockPwSvc, mockAuthSvc)
	ctx := context.Background()
	// Test with empty email
	user1 := &domain.User{Email: "", Password: "password123"}

	result1, err1 := usecase.Create(ctx, user1)
	assert.Error(t, err1)
	assert.Equal(t, domain.ErrInvalidPayload, err1.Error())
	assert.Nil(t, result1)

	// Test with empty password
	user2 := &domain.User{Email: "test@example.com", Password: ""}
	result2, err2 := usecase.Create(ctx, user2)
	assert.Error(t, err2)
	assert.Equal(t, domain.ErrInvalidPayload, err2.Error())
	assert.Nil(t, result2)

	mockRepo.AssertNotCalled(t, "GetUserByEmail", mock.Anything)
	mockPwSvc.AssertNotCalled(t, "Hash", mock.Anything)
}

func TestUserUsecase_Login_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPwSvc := new(mocks.MockPasswordService)
	mockAuthSvc := new(mocks.MockAuthService)

	usecase := NewUserUsecase(mockRepo, mockPwSvc, mockAuthSvc)
	ctx := context.Background()
	userInput := &domain.User{Email: "test@example.com", Password: "password123"}
	storedUser := &domain.User{ID: "1", Email: "test@example.com", Password: "hashedpassword"}
	token := "sometoken"

	mockRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(storedUser, nil)
	mockPwSvc.On("Verify", "hashedpassword", "password123").Return(true, nil)
	mockAuthSvc.On("GenerateToken", "1").Return(token, nil)

	result, tokenResult, err := usecase.Login(ctx, userInput)

	assert.NoError(t, err)
	assert.Equal(t, storedUser, result)
	assert.Equal(t, token, tokenResult)

	mockRepo.AssertExpectations(t)
	mockPwSvc.AssertExpectations(t)
	mockAuthSvc.AssertExpectations(t)
}

func TestUserUsecase_Login_WrongPassword(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPwSvc := new(mocks.MockPasswordService)
	mockAuthSvc := new(mocks.MockAuthService)

	usecase := NewUserUsecase(mockRepo, mockPwSvc, mockAuthSvc)
	ctx := context.Background()
	userInput := &domain.User{Email: "test@example.com", Password: "wrongpassword"}
	storedUser := &domain.User{ID: "1", Email: "test@example.com", Password: "hashedpassword"}

	mockRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(storedUser, nil)
	mockPwSvc.On("Verify", "hashedpassword", "wrongpassword").Return(false, nil)
	result, token, err := usecase.Login(ctx, userInput)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidPayload, err.Error())
	assert.Nil(t, result)
	assert.Empty(t, token)

	mockRepo.AssertExpectations(t)
	mockPwSvc.AssertExpectations(t)
	mockAuthSvc.AssertNotCalled(t, "GenerateToken", mock.Anything)
}

func TestUserUsecase_Login_UserNotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPwSvc := new(mocks.MockPasswordService)
	mockAuthSvc := new(mocks.MockAuthService)

	usecase := NewUserUsecase(mockRepo, mockPwSvc, mockAuthSvc)
	ctx := context.Background()
	userInput := &domain.User{Email: "test@example.com", Password: "password123"}

	mockRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return((*domain.User)(nil), errors.New("not found"))
	result, token, err := usecase.Login(ctx, userInput)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrUserNotFound, err.Error())
	assert.Nil(t, result)
	assert.Empty(t, token)

	mockRepo.AssertExpectations(t)
	mockPwSvc.AssertNotCalled(t, "Verify", mock.Anything)
	mockAuthSvc.AssertNotCalled(t, "GenerateToken", mock.Anything)
}

func TestUserUsecase_Logout(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPwSvc := new(mocks.MockPasswordService)
	mockAuthSvc := new(mocks.MockAuthService)

	usecase := NewUserUsecase(mockRepo, mockPwSvc, mockAuthSvc)

	ctx := context.Background()

	result, err := usecase.Logout(ctx)

	assert.NoError(t, err)
	assert.Equal(t, "logout successful", result)

	// No mocks should be called
	mockRepo.AssertNotCalled(t, mock.Anything)
	mockPwSvc.AssertNotCalled(t, mock.Anything)
	mockAuthSvc.AssertNotCalled(t, mock.Anything)
}

func TestUserUsecase_DeleteUserAndTasks_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPwSvc := new(mocks.MockPasswordService)
	mockAuthSvc := new(mocks.MockAuthService)

	usecase := NewUserUsecase(mockRepo, mockPwSvc, mockAuthSvc)

	ctx := context.Background()
	userID := "1"

	mockRepo.On("Delete", mock.Anything, userID).Return(nil)

	err := usecase.DeleteUserAndTasks(ctx, userID)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestUserUsecase_DeleteUserAndTasks_RepoError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPwSvc := new(mocks.MockPasswordService)
	mockAuthSvc := new(mocks.MockAuthService)

	usecase := NewUserUsecase(mockRepo, mockPwSvc, mockAuthSvc)

	ctx := context.Background()
	userID := "1"

	mockRepo.On("Delete", mock.Anything, userID).Return(errors.New("db error"))

	err := usecase.DeleteUserAndTasks(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInternalServer, err.Error())

	mockRepo.AssertExpectations(t)
}

func TestUserUsecase_VerifyPassword_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPwSvc := new(mocks.MockPasswordService)
	mockAuthSvc := new(mocks.MockAuthService)

	usecase := NewUserUsecase(mockRepo, mockPwSvc, mockAuthSvc)

	ctx := context.Background()
	userID := "1"
	password := "password123"

	storedUser := &domain.User{ID: userID, Email: "test@example.com", Password: "hashedpassword"}

	mockRepo.On("GetUserByID", mock.Anything, userID).Return(storedUser, nil)
	mockPwSvc.On("Verify", "hashedpassword", password).Return(true, nil)

	err := usecase.VerifyPassword(ctx, userID, password)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockPwSvc.AssertExpectations(t)
}

func TestUserUsecase_VerifyPassword_WrongPassword(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPwSvc := new(mocks.MockPasswordService)
	mockAuthSvc := new(mocks.MockAuthService)

	usecase := NewUserUsecase(mockRepo, mockPwSvc, mockAuthSvc)

	ctx := context.Background()
	userID := "1"
	password := "wrongpassword"

	storedUser := &domain.User{ID: userID, Email: "test@example.com", Password: "hashedpassword"}

	mockRepo.On("GetUserByID", mock.Anything, userID).Return(storedUser, nil)
	mockPwSvc.On("Verify", "hashedpassword", password).Return(false, nil)

	err := usecase.VerifyPassword(ctx, userID, password)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidPayload, err.Error())

	mockRepo.AssertExpectations(t)
	mockPwSvc.AssertExpectations(t)
}

func TestUserUsecase_VerifyPassword_UserNotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockPwSvc := new(mocks.MockPasswordService)
	mockAuthSvc := new(mocks.MockAuthService)

	usecase := NewUserUsecase(mockRepo, mockPwSvc, mockAuthSvc)

	ctx := context.Background()
	userID := "1"
	password := "password123"

	mockRepo.On("GetUserByID", mock.Anything, userID).Return((*domain.User)(nil), errors.New("not found"))

	err := usecase.VerifyPassword(ctx, userID, password)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrUserNotFound, err.Error())

	mockRepo.AssertExpectations(t)
	mockPwSvc.AssertNotCalled(t, "Verify", mock.Anything)
}