package pg

import (
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewDatabase(logger *slog.Logger) (*sql.DB, error) {
	connectionString := os.Getenv("DB_CONNECTION_STRING")
	logger.Debug("connecting to database", "connection_string", connectionString)
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
