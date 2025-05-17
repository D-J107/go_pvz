package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDb(ctx context.Context) *DB {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		panic("environment variable DATABASE URL not set")
	}
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		panic("cant establish connection to remote Postgre db")
	}
	return &DB{Pool: pool}
}

func (db *DB) InitDb(ctx context.Context) error {
	creation := []string{`
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		role TEXT NOT NULL CHECK (role IN ('moderator', 'client', 'employee'))
	)`, `
	CREATE TABLE IF NOT EXISTS pvz (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		registration_date_time Timestamp NOT NULL DEFAULT NOW(),
		city TEXT NOT NULL CHECK (city IN ('Москва', 'Санкт-Петербург', 'Казань'))
	) `, `
	CREATE TABLE IF NOT EXISTS receptions (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		date_time TIMESTAMP NOT NULL DEFAULT NOW(),
		pvz_id UUID NOT NULL,
		status TEXT NOT NULL CHECK (status IN ('in_progress', 'close')),
		FOREIGN KEY (pvz_id) REFERENCES pvz(id)
	)`, `
	CREATE TABLE IF NOT EXISTS products (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		date_time TIMESTAMP NOT NULL DEFAULT NOW(),
		type TEXT NOT NULL CHECK (type IN ('электроника', 'одежда', 'обувь')),
		reception_id UUID NOT NULL,
		FOREIGN KEY (reception_id) REFERENCES receptions(id)
	)`}
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		panic("cant open transactions: " + err.Error())
	}
	defer tx.Rollback(ctx)
	for _, createCommand := range creation {
		_, err := tx.Exec(ctx, createCommand)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
