package middleware

import (
	"net/http"

	"github.com/Anurag-S1ngh/attendance-tracker/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthMiddlewareHandler struct {
	authMiddlewareService *service.AuthMiddlewareService
}

func NewMiddlewareHandler(s *service.AuthMiddlewareService) *AuthMiddlewareHandler {
	return &AuthMiddlewareHandler{
		authMiddlewareService: s,
	}
}

func (h *AuthMiddlewareHandler) AuthMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	userID, err := h.authMiddlewareService.AuthMiddleware(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.Set("userID", userID)

	c.Next()
}
