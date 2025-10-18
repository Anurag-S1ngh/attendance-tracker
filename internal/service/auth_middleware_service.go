package service

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

type AuthMiddlewareService struct {
	store *sessions.CookieStore
}

func NewAuthMiddlewareService(store *sessions.CookieStore) *AuthMiddlewareService {
	return &AuthMiddlewareService{
		store: store,
	}
}

func (s *AuthMiddlewareService) AuthMiddleware(ctx *gin.Context) error {
	session, err := s.store.Get(ctx.Request, "user_session")
	if err != nil {
		log.Printf("failed to get user session: %v", err)
		return errors.New("failed to get user session")
	}

	userID, ok := session.Values["user_id"]
	if !ok {
		log.Printf("unauthorized: user_id not found in session")
		return errors.New("unauthorized")
	}

	ctx.Set("userID", userID)

	return nil
}