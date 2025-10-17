package service

import (
	"database/sql"
	"errors"

	db "github.com/Anurag-S1ngh/attendance-tracker/internal/db/generated"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

type AuthService struct {
	db     *db.Queries
	store  *sessions.CookieStore
	logger *zap.Logger
}

func NewAuthService(db *db.Queries, store *sessions.CookieStore, logger *zap.Logger) *AuthService {
	return &AuthService{
		db:     db,
		store:  store,
		logger: logger,
	}
}

func (s *AuthService) SaveSession(ctx *gin.Context, email string) error {
	userUUID, err := s.db.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			userUUID, err = s.db.CreateUser(ctx, db.CreateUserParams{
				ID:    uuid.New(),
				Email: email,
			})
			if err != nil {
				s.logger.Error("failed to get user by email", zap.Error(err))
				return errors.New("internal server error")
			}
		} else {
			return errors.New("internal server error")
		}
	}

	session, err := s.store.Get(ctx.Request, "user_session")
	if err != nil {
		s.logger.Error("failed to get session", zap.Error(err))
		return errors.New("failed to get session")
	}
	session.Values["user_id"] = userUUID.String()
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 day
		HttpOnly: true,
		// TODO: should be true in prod
		Secure: false,
	}

	err = session.Save(ctx.Request, ctx.Writer)
	if err != nil {
		s.logger.Error("failed to save session", zap.Error(err))
		return errors.New("failed to save session")
	}

	return nil
}
