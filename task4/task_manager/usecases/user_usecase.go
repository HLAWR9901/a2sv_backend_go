package usecases

import (
	"context"
	"errors"
	"task_manager/domain"
	"task_manager/infrastructure/config"

	"github.com/google/uuid"
)

type UserUsecase struct {
	repo    domain.IUserRepository
	pwSvc   domain.IPasswordService
	authSvc domain.IAuthService
}

// Init User Usecase
func NewUserUsecase(
	repo domain.IUserRepository,
	pwSvc domain.IPasswordService,
	authSvc domain.IAuthService,
) *UserUsecase {
	return &UserUsecase{
		repo:    repo,
		pwSvc:   pwSvc,
		authSvc: authSvc,
	}
}

// User Usecase application logic
func (uu *UserUsecase) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, config.UserCreationTimeout)
	defer cancel()

	if user.Email == "" || user.Password == "" {
		return nil, errors.New(domain.ErrInvalidPayload)
	}

	existing, _ := uu.repo.GetUserByEmail(ctx, user.Email)
	if existing != nil {
		return nil, errors.New(domain.ErrUserExists)
	}

	hashed, err := uu.pwSvc.Hash(user.Password)
	if err != nil {
		return nil, errors.New(domain.ErrHashingPassword)
	}

	user.ID = uuid.NewString()
	user.Password = hashed

	if err := uu.repo.Create(ctx, user); err != nil {
		return nil, errors.New(domain.ErrInternalServer)
	}

	return user, nil
}

func (uu *UserUsecase) Login(ctx context.Context, user *domain.User) (*domain.User, string, error) {
	ctx, cancel := context.WithTimeout(ctx, config.UserLoginTimeout)
	defer cancel()

	result, err := uu.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return nil, "", errors.New(domain.ErrUserNotFound)
	}

	valid, err := uu.pwSvc.Verify(result.Password, user.Password)
	if err != nil || !valid {
		return nil, "", errors.New(domain.ErrInvalidPayload)
	}

	token, err := uu.authSvc.GenerateToken(result.ID)
	if err != nil {
		return nil, "", errors.New(domain.ErrInternalServer)
	}

	return result, token, nil
}

func (uu *UserUsecase) Logout(ctx context.Context) (string, error) {
	return "logout successful", nil
}

func (uu *UserUsecase) DeleteUserAndTasks(ctx context.Context, userID string) error {
	if err := uu.repo.Delete(ctx, userID); err != nil {
		return errors.New(domain.ErrInternalServer)
	}
	return nil
}

func (uu *UserUsecase) VerifyPassword(ctx context.Context, userID, password string) error {
	user, err := uu.repo.GetUserByID(ctx, userID)
	if err != nil {
		return errors.New(domain.ErrUserNotFound)
	}

	valid, err := uu.pwSvc.Verify(user.Password, password)
	if err != nil || !valid {
		return errors.New(domain.ErrInvalidPayload)
	}

	return nil
}