package cmd

import (
	"context"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
	"gavana/graft"
)

func generateMigrations(models []any) {

	db, err := graft.New()
	if err != nil {
		log.Fatal(err)
	}

	var upQueries []string
	var downQueries []string
	for name, model := range models {
		upQuery,downQuery := generateQueries(strconv.Itoa(name), model)
		upQueries = append(upQueries, upQuery)
		downQueries = append(downQueries, downQuery)
	}

	writeMigrationFile(upQueries, downQueries)
}

func generateQueries(itoa string, model any) (interface{}, interface{}) {

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



