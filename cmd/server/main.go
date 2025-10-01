package main

import (
	"database/sql"
	"log"

	"github.com/Anurag-S1ngh/attendance-tracker/internal/config"
	db "github.com/Anurag-S1ngh/attendance-tracker/internal/db/generated"
	"github.com/Anurag-S1ngh/attendance-tracker/internal/http"
	"github.com/Anurag-S1ngh/attendance-tracker/internal/service"
	_ "github.com/lib/pq"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

func main() {
	cfg := config.Load()
	dbConn, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbConn.Close()
	query := db.New(dbConn)

	authService := service.NewAuthService(query, cfg.JWTSecret)
	attendanceService := service.NewAttendanceService(query)
	eventsService := service.NewEventService(query)
	authMiddlewareService := service.NewAuthMiddlewareService(cfg.JWTSecret, cfg.ClientID, cfg.ClientSecret, cfg.CallBackURL)

	goth.UseProviders(
		google.New(cfg.ClientID, cfg.ClientSecret, cfg.CallBackURL),
	)

	router := http.NewRouter(authService, eventsService, attendanceService, authMiddlewareService)

	log.Println("server is running on port", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
