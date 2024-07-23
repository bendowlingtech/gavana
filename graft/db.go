package graft

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

func New() (*Graft, error) {

	db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic("Couldn't connect to DB")
	}

	return &Graft{db}, nil
}

