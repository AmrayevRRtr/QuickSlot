package main

import "QuickSlot/internal/model"

func loadConfig() *model.MySQLConfig {
	return &model.MySQLConfig{
		Username: "root",
		Password: "Password123",
		Host:     "localhost",
		Port:     "3306",
		DBName:   "QuickSlot",
		SSLMode:  "false",
	}
}
