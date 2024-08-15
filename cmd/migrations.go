package cmd

import (
	"context"
	"fmt"
	"github.com/bendowlingtech/gavana/graft"
	"github.com/spf13/cobra"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
)


func makeMigrationsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "make:migrations",
		Short: "Generate database migrations",
		Run: func(cmd *cobra.Command, args []string) {
			models := []any{

			}
			generateMigrations(models)
		},
	}
}

func generateMigrations(models []any) {
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

func generateQueries(model any, db *graft.Graft) (string, string) {
	rtype := reflect.TypeOf(model)
	mName := rtype.Name()

	if tableExists(db, mName) {
		return generateAlterTableQueries(mName, model)
	} else {
		return generateCreateTableQueries(mName, model)
	}
}

func generateAlterTableQueries(name string, model any) (string, string) {
	var upQuery strings.Builder
	var downQuery strings.Builder

	upQuery.WriteString(fmt.Sprintf("ALTER TABLE %s ", name))
	downQuery.WriteString(fmt.Sprintf("ALTER TABLE %s ", name))

	upQuery.WriteString("ADD COLUMN example_column VARCHAR(255);")
	downQuery.WriteString("DROP COLUMN example_column;")

	return upQuery.String(), downQuery.String()
}

func generateCreateTableQueries(name string, model any) (string, string) {
	var upQuery strings.Builder
	var downQuery strings.Builder

	upQuery.WriteString(fmt.Sprintf("CREATE TABLE %s (", name))
	downQuery.WriteString(fmt.Sprintf("DROP TABLE %s;", name))

	for i := 0; i < reflect.TypeOf(model).NumField(); i++ {
		field := reflect.TypeOf(model).Field(i)
		columnName := field.Name
		columnType := getColumnType(field.Type)

		upQuery.WriteString(fmt.Sprintf("%s %s,", columnName, columnType))
	}

	upQuery.Truncate(upQuery.Len() - 1)
	upQuery.WriteString(");")

	return upQuery.String(), downQuery.String()
}

func getColumnType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64:
		return "INTEGER"
	case reflect.String:
		return "VARCHAR(255)"
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.Float32, reflect.Float64:
		return "REAL"
	default:
		return "TEXT"
	}
}

func tableExists(db *graft.Graft, tableName string) bool {
	var exists bool
	err := db.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)", tableName).Scan(&exists)
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



