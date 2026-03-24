package model

import "time"

type MySQLConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string

	ExecTimeout time.Duration
}
