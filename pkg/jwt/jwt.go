package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Manager struct {
	secret string
}

func NewManager(secret string) *Manager {
	return &Manager{
		secret: secret,
	}
}

func (m *Manager) Generate(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"iat":     time.Now(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secret))
}

func (m *Manager) Verify(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return []byte(m.secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
