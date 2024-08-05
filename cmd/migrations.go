package cmd

import (
	"context"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
	"github.com/bendowlingtech/gavana/graft"
	"github.com/spf13/cobra"
)


func makeMigrationsCmd() *cobra.Command {
	return &cobra.Command{
		Use: "make:migrations",
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
		upQuery,downQuery := generateQueries(model)
		upQueries = append(upQueries, upQuery)
		downQueries = append(downQueries, downQuery)
	}

	writeMigrationFile(upQueries, downQueries)
}

func generateQueries(model any) (string,string) {
	rtype := reflect.TypeOf(model)
	for i := 0; i < rtype.NumField(); i++ {
		name := rtype.Name()
		field := rtype.Field(i)

	}
	return upQuery,downQuery


}

func tableExists(db *graft.Graft, tableName string) bool {
	var exists bool
	err := db.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1", tableName).Scan(exists)
	if err != nil {

	}
	return exists
}

func writeMigrationFile(upQueries, downQueries []string) {
	timestamp := time.Now()


}



