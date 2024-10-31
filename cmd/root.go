/*
Copyright © 2024 Denys <https://github.com/AnyoneClown>
This is my license. There are many like it, but this one is mine.
My license is my best friend. It is my life. I must master it as I must
master my life.
*/
package cmd

import (
	"os"

	"github.com/AnyoneClown/anydb/cmd/backup"
	"github.com/AnyoneClown/anydb/cmd/configure"
	"github.com/AnyoneClown/anydb/cmd/table"
	"github.com/AnyoneClown/anydb/config"
	"github.com/AnyoneClown/anydb/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
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
	utils.InitLogger() // Register Zap logger
	defer utils.Log.Sync()

	utils.CreateFileAndDir()

	var err error
	config.Configs, err = utils.LoadConfigs(config.ConfigFile)
	if err != nil {
		utils.Log.Error("File doesn't exist, creating configuration file", zap.Error(err))
	}

	rootCmd.AddCommand(configure.ConfigureCmd)
	rootCmd.AddCommand(table.TableCmd)
	rootCmd.AddCommand(backup.BackupCmd)
}
