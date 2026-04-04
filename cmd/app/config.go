package main

import (
	"QuickSlot/internal/model"
	"os"

	"github.com/joho/godotenv"
)

func loadConfig() *model.MySQLConfig {
	_ = godotenv.Load()

	return &model.MySQLConfig{
		Username: getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASS", "Password123"),
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3306"),
		DBName:   getEnv("DB_NAME", "QuickSlot"),
		SSLMode:  getEnv("DB_SSL", "false"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
