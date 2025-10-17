package middleware

import (
	"fmt"
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
	err := h.authMiddlewareService.AuthMiddleware(c)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.Next()
}
