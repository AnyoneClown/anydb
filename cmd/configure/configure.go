/*
Copyright © 2024 Denys <https://github.com/AnyoneClown>
This is my license. There are many like it, but this one is mine.
My license is my best friend. It is my life. I must master it as I must
master my life.
*/
package configure

import (
	"fmt"
	"os"

	"github.com/AnyoneClown/anydb/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var ConfigureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure your database credentials",
	Long:  `Use it to choose database credentials. You can add, edit, remove, and list your configurations!`,
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(newModel("Database configuration"), tea.WithAltScreen())

		finalModel, err := p.Run()
		if err != nil {
			fmt.Printf("Error running program: %v\n", err)
			return
		}

		if finalModel.(model).choice != nil {
			choice := finalModel.(model).choice
			fmt.Printf("Selected configuration: %s\n", choice.ConfigName)

			data, err := yaml.Marshal(choice)
			if err != nil {
				return
			}
			os.WriteFile(utils.DefaultConfigFile, data, 0644)
		}
	},
}

func init() {
	utils.CreateFileAndDir()

	var err error
	utils.Configs, err = utils.LoadConfigs(utils.ConfigFile)
	if err != nil {
		fmt.Println("File doesn't exist, creating configuration file")
	}

	ConfigureCmd.AddCommand(addCmd)
	ConfigureCmd.AddCommand(removeCmd)
}
