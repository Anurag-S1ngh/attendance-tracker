package main

import (
	"database/sql"
	"encoding/hex"
	"log"

	"github.com/Anurag-S1ngh/attendance-tracker/internal/config"
	db "github.com/Anurag-S1ngh/attendance-tracker/internal/db/generated"
	"github.com/Anurag-S1ngh/attendance-tracker/internal/http"
	"github.com/Anurag-S1ngh/attendance-tracker/internal/service"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

var Store *sessions.CookieStore

func main() {
	cfg := config.Load()
	dbConn, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbConn.Close()
	query := db.New(dbConn)

	hashKey, err := hex.DecodeString(cfg.HASHKEYHEX)
	if err != nil {
		log.Fatal(err)
	}

	blockKey, err := hex.DecodeString(cfg.BLOCKKEYHEX)
	if err != nil {
		log.Fatal(err)
	}

	Store = sessions.NewCookieStore(hashKey, blockKey)

	authService := service.NewAuthService(query, Store)
	attendanceService := service.NewAttendanceService(query)
	eventsService := service.NewEventService(query)
	authMiddlewareService := service.NewAuthMiddlewareService(Store)

	goth.UseProviders(
		google.New(cfg.ClientID, cfg.ClientSecret, cfg.CallBackURL),
	)
	gothic.Store = Store

	router := http.NewRouter(authService, eventsService, attendanceService, authMiddlewareService)

	log.Printf("server is running on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
