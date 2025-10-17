package http

import (
	"time"

	"github.com/Anurag-S1ngh/attendance-tracker/internal/http/handlers"
	"github.com/Anurag-S1ngh/attendance-tracker/internal/http/middleware"
	"github.com/Anurag-S1ngh/attendance-tracker/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(authService *service.AuthService, eventsService *service.EventsService, attendanceService *service.AttendanceService, authMiddlwareService *service.AuthMiddlewareService) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	authHandler := handlers.NewAuthHandler(authService)
	healthCheckHandler := handlers.NewHealthCheckHandler()
	eventsHandler := handlers.NewEventsHandler(eventsService)
	attendanceHandler := handlers.NewAttendanceHandler(attendanceService)
	authMiddlewareHandler := middleware.NewMiddlewareHandler(authMiddlwareService)

	r.GET("/health_check", healthCheckHandler.HealthCheck)

	v1 := r.Group("/api/v1")
	{
		r.GET("/auth/:provider", authHandler.SignInWithProvider)
		r.GET("/auth/:provider/callback", authHandler.CallbackHandler)

		v1.GET("/event", authMiddlewareHandler.AuthMiddleware, eventsHandler.GetAllEvents)
		v1.GET("/event/:eventID", authMiddlewareHandler.AuthMiddleware, eventsHandler.GetEvent)
		v1.POST("/event", authMiddlewareHandler.AuthMiddleware, eventsHandler.CreateEvent)
		v1.DELETE("/event/:eventID", authMiddlewareHandler.AuthMiddleware, eventsHandler.DeleteEvent)

		v1.GET("/attendance", authMiddlewareHandler.AuthMiddleware, attendanceHandler.GetAttendance)
		v1.POST("/attendance/:eventID", authMiddlewareHandler.AuthMiddleware, attendanceHandler.MarkAttendance)
		v1.DELETE("/attendance/:attendanceID", authMiddlewareHandler.AuthMiddleware, attendanceHandler.DeleteAttendance)
	}

	return r
}
