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
	"go.uber.org/zap"
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

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	authService := service.NewAuthService(query, Store, logger)
	attendanceService := service.NewAttendanceService(query, logger)
	eventsService := service.NewEventService(query, logger)
	authMiddlewareService := service.NewAuthMiddlewareService(Store, logger)

	goth.UseProviders(
		google.New(cfg.ClientID, cfg.ClientSecret, cfg.CallBackURL),
	)
	gothic.Store = Store

	router := http.NewRouter(authService, eventsService, attendanceService, authMiddlewareService)

	logger.Info("server is running on port", zap.String("port", cfg.Port))
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
