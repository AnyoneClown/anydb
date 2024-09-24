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

			err := utils.LoadDefaultConfig()
			if err != nil {
				return
			}

			if *choice == utils.DefaultConfigData {
				err = os.Truncate(utils.DefaultConfigFile, 0)
				if err != nil {
					fmt.Printf("Error truncating file: %v\n", err)
					return
				}
			}

			for index, value := range utils.Configs {
				if value == *choice {
					utils.Configs = append(utils.Configs[:index], utils.Configs[index+1:]...)
				}
			}

			if err := utils.SaveConfigs(utils.Configs, utils.ConfigFile); err != nil {
				fmt.Printf("Failed to save configuration: %v\n", err)
			}
		}
	},
}

func init() {
	removeCmd.Flags().BoolP("help", "h", false, "help for add")
	removeCmd.Flags().MarkHidden("help")
}
