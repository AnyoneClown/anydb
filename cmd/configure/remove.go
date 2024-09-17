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

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove database configuration",
	Run: func(cmd *cobra.Command, args []string) {

		p := tea.NewProgram(newModel("Remove configuration"), tea.WithAltScreen())

		finalModel, err := p.Run()
		if err != nil {
			fmt.Printf("Error running program: %v\n", err)
			return
		}

		if finalModel.(model).choice != nil {
			choice := finalModel.(model).choice
			fmt.Printf("Deleted configuration: %s\n", choice.ConfigName)

			err := loadDefaultConfig()
			if err != nil {
				return
			}

			if *choice == defaultConfigData {
				err = os.Truncate(defaultConfigFile, 0)
				if err != nil {
					fmt.Printf("Error truncating file: %v\n", err)
					return
				}
			}

			for index, value := range configs {
				if value == *choice {
					configs = append(configs[:index], configs[index+1:]...)
				}
			}

			if err := saveConfigs(configs, configFile); err != nil {
				fmt.Printf("Failed to save configuration: %v\n", err)
			}
		}
	},
}

func init() {
	removeCmd.Flags().BoolP("help", "h", false, "help for add")
	removeCmd.Flags().MarkHidden("help")
}
