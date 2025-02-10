package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
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
