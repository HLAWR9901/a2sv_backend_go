package passwordservice

import (
	"golang.org/x/crypto/bcrypt"
	"task_manager/domain"
)

const DefaultCost = bcrypt.DefaultCost

type BcryptService struct {
	Cost int
}

func NewBcryptService(cost int) domain.IPasswordService {
	return &BcryptService{Cost: cost}
}

func (s *BcryptService) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.Cost)
	return string(bytes), err
}

func (s *BcryptService) Verify(hash, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}
	return err == nil, err
}