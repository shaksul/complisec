package database

import (
	"database/sql"
	"fmt"

	"risknexus/backend/internal/repo"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(databaseURL string) (*repo.DB, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &repo.DB{db}, nil
}
