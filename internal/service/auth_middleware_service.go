package service

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

type AuthMiddlewareService struct {
	store  *sessions.CookieStore
	logger *zap.Logger
}

func NewAuthMiddlewareService(store *sessions.CookieStore, logger *zap.Logger) *AuthMiddlewareService {
	return &AuthMiddlewareService{
		store:  store,
		logger: logger,
	}
}

func (s *AuthMiddlewareService) AuthMiddleware(ctx *gin.Context) error {
	session, err := s.store.Get(ctx.Request, "user_session")
	if err != nil {
		s.logger.Error("failed to get user session", zap.Error(err))
		return errors.New("failed to get user session")
	}

	userID, ok := session.Values["user_id"]
	if !ok {
		s.logger.Error("unauthorized", zap.Error(err))
		return errors.New("unauthorized")
	}

	ctx.Set("userID", userID)

	return nil
}
