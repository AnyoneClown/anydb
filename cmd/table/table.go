/*
Copyright © 2024 Denys <https://github.com/AnyoneClown>
This is my license. There are many like it, but this one is mine.
My license is my best friend. It is my life. I must master it as I must
master my life.
*/
package table

import (
	"fmt"

	"github.com/AnyoneClown/anydb/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	_ "github.com/lib/pq"
)

var TableCmd = &cobra.Command{
	Use:   "table",
	Short: "Display tables and their contents",
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("rows")

		dsn, err := utils.GetDBString()
		if err != nil {
			fmt.Println("Error getting database string:", err)
			return
		}

		db, err := sqlx.Connect("postgres", dsn)
		if err != nil {
			fmt.Println("Error connecting to database:", err)
			return
		}
		defer db.Close()

		model, err := NewModel(db, limit)
		if err != nil {
			fmt.Println("Error initializing model:", err)
			return
		}
		if _, err := tea.NewProgram(model).Run(); err != nil {
			fmt.Println("Error running program:", err)
			return
		}
	},
}

func init() {
	TableCmd.Flags().IntP("rows", "r", 5, "Number of rows to display")
}
