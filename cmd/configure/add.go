/*
Copyright © 2024 Denys <https://github.com/AnyoneClown>
This is my license. There are many like it, but this one is mine.
My license is my best friend. It is my life. I must master it as I must
master my life.

*/
package configure

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type DBConfig struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

var configs []DBConfig
var configFile string

func loadConfigs() {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return
	}
	json.Unmarshal(data, &configs)
}

func saveConfigs() {
	data, _ := json.MarshalIndent(configs, "", "  ")
	os.WriteFile(configFile, data, 0644)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new database configuration",

	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		database, _ := cmd.Flags().GetString("database")

		config := DBConfig{
			Name:     name,
			Host:     host,
			Port:     port,
			User:     user,
			Password: password,
			Database: database,
		}

		configs = append(configs, config)
		saveConfigs()
	},
}

func init() {
	homeDir, _ := os.UserHomeDir()
	configFile = filepath.Join(homeDir, "anydb-config.json")
	loadConfigs()

	addCmd.Flags().String("name", "", "Configuration name")
	addCmd.Flags().String("host", "localhost", "Database host")
	addCmd.Flags().Int("port", 5432, "Database port")
	addCmd.Flags().String("user", "", "Database user")
	addCmd.Flags().String("password", "", "Database password")
	addCmd.Flags().String("database", "", "Database name")

	addCmd.MarkFlagRequired("name")
	addCmd.MarkFlagRequired("host")
	addCmd.MarkFlagRequired("port")
	addCmd.MarkFlagRequired("user")
	addCmd.MarkFlagRequired("password")
	addCmd.MarkFlagRequired("database")

	addCmd.Flags().BoolP("help", "h", false, "help for add")
	addCmd.Flags().MarkHidden("help")
}
