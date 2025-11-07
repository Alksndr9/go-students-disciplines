package psql

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"

	"gitgub.com/Alksndr9/go-students-disciplines/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/database"
	"go.uber.org/zap"
)

func Connect(ctx context.Context, cfg *config.Config, log *zap.Logger) (*pgxpool.Pool, error) {
	connString := (&url.URL{
		Scheme: "postgres",

		User: url.UserPassword(
			cfg.User,
			cfg.Password,
		),
		Host: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Path: cfg.Database,
	}).String()

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if err = pingDB(ctx, db, log); err != nil {
		return nil, err
	}

	stdDB := stdlib.OpenDB(*poolConfig.ConnConfig)
	defer stdDB.Close()

	g, err := goose.NewProvider(
		database.DialectPostgres,
		stdDB,
		os.DirFS("./db/migrations"),
	)

	if err != nil {
		return nil, err
	}

	if _, err = g.Up(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func pingDB(ctx context.Context, db *pgxpool.Pool, log *zap.Logger) error {
	if db == nil {
		return errors.New("db is nill")
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	timeout := time.After(time.Second * 10)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout:
			return errors.New("db ping timeout")
		case <-ticker.C:
			if err := db.Ping(ctx); err == nil {
				return nil
			} else {
				log.Error("ping error")
			}
		}
	}
}
