package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DBURL        string
	ClientID     string
	ClientSecret string
	CallBackURL  string
	HASHKEYHEX   string
	BLOCKKEYHEX  string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return nil
	}
	return &Config{
		Port:         getEnv("PORT", "8000"),
		DBURL:        getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"),
		ClientID:     getEnv("CLIENT_ID", ""),
		ClientSecret: getEnv("CLIENT_SECRET", ""),
		CallBackURL:  getEnv("CLIENT_CALLBACK_URL", "4d0a5003a7ada6c"),
		HASHKEYHEX:   getEnv("HASHKEYHEX", "52007e5b1f584"),
		BLOCKKEYHEX:  getEnv("BLOCKKEYHEX", "52007e5b1f584"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return fallback
}
