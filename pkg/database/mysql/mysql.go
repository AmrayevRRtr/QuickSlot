package mysql

import (
	"QuickSlot/internal/model"
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

type Dialect struct {
	DB *sqlx.DB
}

func NewMySQLDialect(ctx context.Context, cfg *model.MySQLConfig) *Dialect {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&tls=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return &Dialect{DB: db}
}

func AutoMigrate(cfg *model.MySQLConfig) {
	sourceURL := "file://database/migrations"

	databaseURL := fmt.Sprintf(
		"mysql://%s:%s@tcp(%s:%s)/%s?tls=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
}
