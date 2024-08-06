package graft

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
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
	query := qb.buildQuery(true)
	return qb.graft.db.QueryRow(query, qb.args...).Scan(dest)
}

func (qb *QueryBuilder) buildQuery(singleRecord bool) string {
	query := "SELECT "

	if len(qb.selectCols) > 0 {
		query += strings.Join(qb.selectCols, ", ")
	} else {
		query += "*"
	}

	query += " FROM " + qb.table

	if len(qb.whereClauses) > 0 {

	}

	if len(qb.orderByClauses) > 0 {
		query += " ORDER BY " + strings.Join(qb.orderByClauses, ", ")
	}

	if singleRecord {
		query += " LIMIT 1"
	}
	return query
}





