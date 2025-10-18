package service

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	db "github.com/Anurag-S1ngh/attendance-tracker/internal/db/generated"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

type AuthService struct {
	db    *db.Queries
	store *sessions.CookieStore
}

func NewAuthService(db *db.Queries, store *sessions.CookieStore) *AuthService {
	return &AuthService{
		db:    db,
		store: store,
	}
}

func (s *AuthService) SaveSession(ctx *gin.Context, email string) error {
	userUUID, err := s.db.GetUserByEmail(ctx, email)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			softDeletedUserUUID, err := s.db.GetSoftDeletedUserByEmail(ctx, email)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					userUUID, err = s.db.CreateUser(ctx, db.CreateUserParams{
						ID:    uuid.New(),
						Email: email,
					})
					if err != nil {
						log.Printf("failed to create user: %v", err)
						return errors.New("internal server error")
					}
				} else {
					log.Printf("failed to get soft deleted user by email: %v", err)
					return errors.New("internal server error")
				}
			} else {
				userUUID, err = s.db.UndeleteUser(ctx, softDeletedUserUUID)
				if err != nil {
					log.Printf("failed to undelete user: %v", err)
					return errors.New("internal server error")
				}
			}
		} else {
			log.Printf("failed to get user by email: %v", err)
			return errors.New("internal server error")
		}
	}

	session, err := s.store.Get(ctx.Request, "user_session")
	if err != nil {
		log.Printf("failed to get session: %v", err)
		return errors.New("failed to get session")
	}

	session.Values["user_id"] = userUUID.String()
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 day
		HttpOnly: true,
		Secure:   false, // TODO: should be true in prod
	}

	err = session.Save(ctx.Request, ctx.Writer)
	if err != nil {
		log.Printf("failed to save session: %v", err)
		return errors.New("failed to save session")
	}

	return nil
}
