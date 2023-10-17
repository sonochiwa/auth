package db

import (
	"embed"
	"fmt"
	"github.com/pressly/goose/v3"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sonochiwa/auth/internal/config"
)

func NewDBConnection(cfg *config.DBConnectionConfig) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	switch cfg.Type {
	case "postgres":
		db, err = sqlx.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s "+
			"sslmode=disalbe", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database))
	case "mysql":
		db, err = sqlx.Open("mysql", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s "+
			"sslmode=disalbe", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database))
	default:
		return nil, fmt.Errorf("unknown db type %s", cfg.Type)
	}

	if err != nil {
		return nil, err
	}

	db.DB.SetConnMaxLifetime(time.Duration(cfg.Pool.MaxOpenConns))
	db.SetMaxIdleConns(cfg.Pool.MaxIdleConns)
	db.SetMaxIdleConns(cfg.Pool.MaxOpenConns)

	return db, nil
}

var embedMigrations embed.FS

func ApplyMigrations(dbType string, db *sqlx.DB) error {
	goose.SetBaseFS(embedMigrations)

	var dialect string
	switch dbType {
	case "postgres":
		dialect = "postgres"
	case "mysql":
		dialect = "mysql"
	}

	if len(dialect) == 0 {
		return fmt.Errorf("unknown db type %s", dbType)
	}

	if err := goose.SetDialect(dialect); err != nil {
		return err
	}

	if err := goose.Up(db.DB, "migrations"); err != nil {
		return err
	}

	return nil
}
