package auth

import (
	"errors"
	"task_manager/domain"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTService struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTService(secret string, ttl time.Duration) domain.IAuthService {
	return &JWTService{secret: []byte(secret), ttl: ttl}
}

func (s *JWTService) GenerateToken(userID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.ttl)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *JWTService) ValidateToken(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{},
		func(t *jwt.Token) (interface{}, error) { return s.secret, nil })
	
	if err != nil || !token.Valid {
		return "", errors.New(domain.ErrInvalidToken)
	}
	
	claims := token.Claims.(*jwt.RegisteredClaims)
	return claims.Subject, nil
}