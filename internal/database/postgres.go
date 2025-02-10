package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
)

type Postgres struct {
	Host, Port, User, Password, Name string
}

func NewPostgresDB(options Postgres) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			options.User,
			options.Password,
			options.Host,
			options.Port,
			options.Name))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func PostgresInit() *Postgres {
	return &Postgres{
		Name:     os.Getenv("POSTGRES_DBNAME"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Host:     os.Getenv("POSTGRES_HOST"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	}
}
