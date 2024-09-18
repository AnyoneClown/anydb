/*
Copyright © 2024 Denys <https://github.com/AnyoneClown>
This is my license. There are many like it, but this one is mine.
My license is my best friend. It is my life. I must master it as I must
master my life.
*/
package last

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var LastCmd = &cobra.Command{
	Use:   "last [table_name]",
	Short: "Display the last 5 records from the specified table",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tableName := args[0]
		rows, _ := cmd.Flags().GetInt("rows")

		dsn := "postgresql://postgres:AA6tgixbKoCkb61S0cUtuA@anyone-9715.7tc.aws-eu-central-1.cockroachlabs.cloud:26257/CocaCallsAPI?sslmode=verify-full"

		db, err := sqlx.Connect("postgres", dsn)
		if err != nil {
			return
		}
		defer db.Close()

		model, err := initializeModel(db, tableName, rows)
		if err != nil {
			fmt.Println("Error initializing model:", err)
			os.Exit(1)
		}
		if _, err := tea.NewProgram(model).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	},
}

func init() {
	LastCmd.Flags().IntP("rows", "r", 5, "number of rows to collect from table[Default 5]")
}