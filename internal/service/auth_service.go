package service

import (
	"context"
	"errors"

	db "github.com/Anurag-S1ngh/attendance-tracker/internal/db/generated"
	"github.com/Anurag-S1ngh/attendance-tracker/pkg/bcrypt"
	"github.com/Anurag-S1ngh/attendance-tracker/pkg/jwt"
	"github.com/google/uuid"
)

type AuthService struct {
	db  *db.Queries
	jwt *jwt.Manager
}

func NewAuthService(db *db.Queries, secret string) *AuthService {
	return &AuthService{
		db:  db,
		jwt: jwt.NewManager(secret),
	}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.db.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("user not found")
	}

	if !bcrypt.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("invalid credentials")
	}

	return s.jwt.Generate(user.ID)
}

func (s *AuthService) Register(ctx context.Context, email, password string) error {
	_, err := s.db.GetUserByEmail(ctx, email)
	if err == nil {
		return errors.New("user already exists")
	}
	hashPassword, err := bcrypt.HashPassword(password)
	if err != nil {
		return err
	}
	return s.db.CreateUser(ctx, db.CreateUserParams{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hashPassword,
	})
}
