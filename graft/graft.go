package graft

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Graft struct {
	*pgxpool.Pool
}


