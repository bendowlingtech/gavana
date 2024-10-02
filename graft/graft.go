package graft

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Graft struct {
	Db *pgxpool.Pool
}

type QueryBuilder struct {
	graft          *Graft
	table          string
	selectCols     []string
	whereClauses   []string
	orderByClauses []string
	limitClause    int
	args           []interface{}
	argCounter     int
}

func (g *Graft) Table(tableName string) *QueryBuilder {
	return &QueryBuilder{
		graft:       g,
		table:       tableName,
		limitClause: -1,
		argCounter:  1,
	}
}

func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
	qb.selectCols = columns
	return qb
}

func (qb *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
	for _, arg := range args {
		placeholder := fmt.Sprintf("$%d", qb.argCounter)
		condition = strings.Replace(condition, "?", placeholder, 1)
		qb.argCounter++
		qb.args = append(qb.args, arg)
	}
	qb.whereClauses = append(qb.whereClauses, condition)
	return qb
}

func (qb *QueryBuilder) OrderBy(clause string) *QueryBuilder {
	qb.orderByClauses = append(qb.orderByClauses, clause)
	return qb
}

func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limitClause = limit
	return qb
}

func (qb *QueryBuilder) First(ctx context.Context, dest ...interface{}) error {
	query := qb.buildQuery(true)
	row := qb.graft.Db.QueryRow(ctx, query, qb.args...)
	return row.Scan(dest...)
}

func (qb *QueryBuilder) All(ctx context.Context, dest interface{}) error {
	rows, err := qb.graft.Db.Query(ctx, qb.buildQuery(false), qb.args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	result, ok := dest.(*[]map[string]interface{})
	if !ok {
		return fmt.Errorf("dest must be of type *[]map[string]interface{}")
	}

	fieldDescriptions := rows.FieldDescriptions()
	columns := make([]string, len(fieldDescriptions))
	for i, fd := range fieldDescriptions {
		columns[i] = string(fd.Name)
	}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return err
		}

		rowMap := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				rowMap[col] = string(b)
			} else {
				rowMap[col] = val
			}
		}

		*result = append(*result, rowMap)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func (qb *QueryBuilder) buildQuery(singleRecord bool) string {
	var query strings.Builder

	query.WriteString("SELECT ")
	if len(qb.selectCols) > 0 {
		query.WriteString(strings.Join(qb.selectCols, ", "))
	} else {
		query.WriteString("*")
	}

	query.WriteString(" FROM ")
	query.WriteString(qb.table)

	if len(qb.whereClauses) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(qb.whereClauses, " AND "))
	}

	if len(qb.orderByClauses) > 0 {
		query.WriteString(" ORDER BY ")
		query.WriteString(strings.Join(qb.orderByClauses, ", "))
	}

	if singleRecord && qb.limitClause == -1 {
		query.WriteString(" LIMIT 1")
	} else if qb.limitClause > 0 {
		query.WriteString(fmt.Sprintf(" LIMIT %d", qb.limitClause))
	}

	return query.String()
}
