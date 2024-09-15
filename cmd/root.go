/*
Copyright © 2024 Denys <https://github.com/AnyoneClown>
This is my license. There are many like it, but this one is mine.
My license is my best friend. It is my life. I must master it as I must
master my life.

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/AnyoneClown/anydb/cmd/configure"
)

var rootCmd = &cobra.Command{
	Use:   "anydb",
	Short: "CLI tool for managing your DB. Get table content, backup your DB!",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(configure.ConfigureCmd)
}
