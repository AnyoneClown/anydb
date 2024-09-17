package configure

import (
	"fmt"
	"os"

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
			os.WriteFile(defaultConfigFile, data, 0644)
		}
	},
}

func init() {
	createFileAndDir()

	var err error
	configs, err = loadConfigs(configFile)
	if err != nil {
		fmt.Println("File doesn't exist, creating configuration file")
	}

	ConfigureCmd.AddCommand(addCmd)
	ConfigureCmd.AddCommand(removeCmd)
}
