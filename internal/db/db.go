package db

import (
	"database/sql"
	"fmt"
	"time"

	// Imported for using postgres driver.
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/ivanterekh/go-skeleton/internal/env"
)

var cfg = struct {
	user, password, address, name, sslMode string
	pingRetries                            int
	pingInterval                           time.Duration
}{
	user:     env.GetString("DB_USER", "postgres"),
	password: env.GetString("DB_PASSWORD", "password"),
	address:  env.GetString("DB_ADDRESS", "localhost"),
	name:     env.GetString("DB_NAME", "goskeleton"),
	sslMode:  env.GetString("DB_SSL_MODE", "disable"),

	pingRetries:  env.GetInt("DB_PING_RETRIES", 3),
	pingInterval: env.GetDuration("DB_PING_INTERVAL", time.Second*4),
}

// New creates db instance and pings it.
func New() (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.user,
		cfg.password,
		cfg.address,
		cfg.name,
		cfg.sslMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to database")
	}

	if err := ping(db); err != nil {
		return nil, errors.Wrap(err, "could not ping database")
	}

	return db, nil
}

func ping(db *sql.DB) error {
	var err error
	for i := 0; i < cfg.pingRetries; i++ {
		if err = db.Ping(); err == nil {
			return nil
		}
		time.Sleep(cfg.pingInterval)
	}

	return err
}
