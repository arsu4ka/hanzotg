package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresOptions struct {
	Host     string
	Port     string
	Username string
	Password string
	DbName   string
}

func OpenPostgres(options PostgresOptions) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		options.Host,
		options.Port,
		options.Username,
		options.Password,
		options.DbName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
