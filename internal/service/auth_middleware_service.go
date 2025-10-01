package service

import (
	"errors"

	"github.com/Anurag-S1ngh/attendance-tracker/pkg/jwt"
)

type AuthMiddlewareService struct {
	jwt          *jwt.Manager
	clientID     string
	clientSecret string
	callBackURL  string
}

func NewAuthMiddlewareService(jwtSecret, clientID, clientSecret, callBackURL string) *AuthMiddlewareService {
	return &AuthMiddlewareService{
		jwt:          jwt.NewManager(jwtSecret),
		clientID:     clientID,
		clientSecret: clientSecret,
		callBackURL:  callBackURL,
	}
}

func (s *AuthMiddlewareService) AuthMiddleware(token string) (string, error) {
	if token == "" {
		return "", errors.New("unauthorized")
	}

	decoded, err := s.jwt.Verify(token)
	if err != nil {
		return "", errors.New("unauthorized")
	}

	if decoded["user_id"] == nil {
		return "", errors.New("unauthorized")
	}

	return decoded["user_id"].(string), nil
}
