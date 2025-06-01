package pg

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewDatabase() (*sql.DB, error) {
	connectionString := os.Getenv("DB_CONNECTION_STRING")
	log.Printf("Connecting to database at: %s\n", connectionString)
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
