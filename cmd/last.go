/*
Copyright © 2024 Denys <https://github.com/AnyoneClown>
This is my license. There are many like it, but this one is mine.
My license is my best friend. It is my life. I must master it as I must
master my life.

*/
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func getLastFiveRecords(db *sqlx.DB, tableName string) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id DESC LIMIT 5", tableName)
	rows, err := db.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		result := make(map[string]interface{})
		err := rows.MapScan(result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

var lastCmd = &cobra.Command{
	Use:   "last [table_name]",
	Short: "Display the last 5 records from the specified table",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tableName := args[0]
		
		dsn := "postgresql://postgres:AA6tgixbKoCkb61S0cUtuA@anyone-9715.7tc.aws-eu-central-1.cockroachlabs.cloud:26257/CocaCallsAPI?sslmode=verify-full"

		db, err := sqlx.Connect("postgres", dsn)
		if err != nil {
			return
		}
		defer db.Close()

		records, err := getLastFiveRecords(db, tableName)
		if err != nil {
			log.Fatalf("Failed to get last 5 records: %v", err)
		}

		if len(records) == 0 {
			fmt.Println("No records found.")
			return
		}

		for i, record := range records {
			fmt.Printf("Record %d:\n", i+1)
			for key, value := range record {
				fmt.Printf("  %s: %v\n", key, value)
			}
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(lastCmd)
}