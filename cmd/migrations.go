package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/bendowlingtech/gavana/graft"
	"github.com/spf13/cobra"
)

func makeMigrationsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "make:migrations",
		Short: "Generate database migrations",
		Run: func(cmd *cobra.Command, args []string) {
			var models []interface{}
			generateMigrations(models)
		},
	}
}

func generateMigrations(models []interface{}) {
	db, err := graft.New()
	if err != nil {
		log.Fatal(err)
	}

	var upQueries []string
	var downQueries []string
	for _, model := range models {
		upQuery, downQuery := generateQueries(model, db)
		upQueries = append(upQueries, upQuery)
		downQueries = append(downQueries, downQuery)
	}

	writeMigrationFile(upQueries, downQueries)
}

func generateQueries(model interface{}, db *graft.Graft) (string, string) {
	rtype := reflect.TypeOf(model)
	mName := rtype.Name()

	if tableExists(db, mName) {
		return generateAlterTableQueries(mName, model, db)
	} else {
		return generateCreateTableQueries(mName, model)
	}
}

func generateAlterTableQueries(tableName string, model interface{}, db *graft.Graft) (string, string) {

	return "", ""
}

func generateCreateTableQueries(tableName string, model interface{}) (string, string) {
	var upQuery strings.Builder
	var downQuery strings.Builder

	upQuery.WriteString(fmt.Sprintf("CREATE TABLE %s (\n", tableName))
	downQuery.WriteString(fmt.Sprintf("DROP TABLE %s;", tableName))

	fields := reflect.VisibleFields(reflect.TypeOf(model))
	for i, field := range fields {
		columnDefinition := generateColumnDefinition(field)
		upQuery.WriteString("    " + columnDefinition)
		if i < len(fields)-1 {
			upQuery.WriteString(",\n")
		} else {
			upQuery.WriteString("\n")
		}
	}
	upQuery.WriteString(");")

	return upQuery.String(), downQuery.String()
}

func generateColumnDefinition(field reflect.StructField) string {
	columnName := field.Name
	columnType := getColumnType(field.Type)
	tags := parseGraftTags(field.Tag)

	var columnConstraints []string

	if tags["column"] != "" {
		columnName = tags["column"]
	}

	if tags["type"] != "" {
		columnType = tags["type"]
	}

	if tags["primaryKey"] == "true" {
		columnConstraints = append(columnConstraints, "PRIMARY KEY")
	}

	if tags["unique"] == "true" {
		columnConstraints = append(columnConstraints, "UNIQUE")
	}

	if tags["notNull"] == "true" {
		columnConstraints = append(columnConstraints, "NOT NULL")
	}

	if tags["default"] != "" {
		columnConstraints = append(columnConstraints, fmt.Sprintf("DEFAULT %s", tags["default"]))
	}

	columnDefinition := fmt.Sprintf("%s %s", columnName, columnType)
	if len(columnConstraints) > 0 {
		columnDefinition += " " + strings.Join(columnConstraints, " ")
	}

	return columnDefinition
}

func parseGraftTags(tag reflect.StructTag) map[string]string {
	tags := make(map[string]string)
	graftTag := tag.Get("graft")
	if graftTag == "" {
		return tags
	}

	tagParts := strings.Split(graftTag, ";")
	for _, part := range tagParts {
		kv := strings.SplitN(part, ":", 2)
		if len(kv) == 2 {
			tags[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		} else if len(kv) == 1 {
			tags[strings.TrimSpace(kv[0])] = "true"
		}
	}
	return tags
}

func getColumnType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64:
		return "INTEGER"
	case reflect.Uint, reflect.Uint32, reflect.Uint64:
		return "INTEGER"
	case reflect.String:
		return "VARCHAR(255)"
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.Float32, reflect.Float64:
		return "REAL"
	case reflect.Struct:
		if t.Name() == "Time" {
			return "TIMESTAMP"
		}
	default:
		return "TEXT"
	}
	return "TEXT"
}

func tableExists(g *graft.Graft, tableName string) bool {
	var exists bool
	err := g.Db.QueryRow(context.Background(),
		"SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)",
		strings.ToLower(tableName)).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}
	return exists
}

func writeMigrationFile(upQueries, downQueries []string) {
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("migrations/%s_migration.sql", timestamp)

	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.WriteString("-- Up Migration\n")
	if err != nil {
		log.Fatal(err)
	}
	for _, query := range upQueries {
		_, err = f.WriteString(query + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err = f.WriteString("\n-- Down Migration\n")
	if err != nil {
		log.Fatal(err)
	}
	for _, query := range downQueries {
		_, err = f.WriteString(query + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}
