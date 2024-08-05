package graft

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Graft struct {
	db *pgxpool.Pool
}

type QueryBuilder struct {
	graft *Graft
	table string
	selectCols []string
	whereClauses []string
	orderByClauses []string
	limitClause []string
	args []interface{}
}

func(g *Graft) Table(tableName string) *QueryBuilder {
	return &QueryBuilder{
		graft: g,
		table: tableName,
	}
}

func (qb QueryBuilder) Select(columns ...string) *QueryBuilder {
	qb.selectCols = columns
	return &qb
}

func (qb *QueryBuilder) First() error{
	query := qb.
}

func (qb *QueryBuilder) buildQuery(singleRecord bool) string {
	query := "SELECT "

	if len(qb.selectCols) > 0 {

	} else {

	}
}





