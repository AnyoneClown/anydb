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
	"strings"

	"github.com/AnyoneClown/anydb/config"
	"github.com/AnyoneClown/anydb/utils"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	errorStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type addModel struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
	errors     []string
}

func initialModel() addModel {
	m := addModel{
		inputs: make([]textinput.Model, 7),
		errors: make([]string, 7),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Config name"
			t.Focus()
			t.Validate = utils.ValidateNotEmpty
		case 1:
			t.Placeholder = "Host"
			t.CharLimit = 64
			t.Validate = utils.ValidateNotEmpty
		case 2:
			t.Placeholder = "Port"
			t.CharLimit = 64
			t.Validate = utils.ValidatePort
		case 3:
			t.Placeholder = "User"
			t.CharLimit = 64
			t.Validate = utils.ValidateNotEmpty
		case 4:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
			t.Validate = utils.ValidateNotEmpty
		case 5:
			t.Placeholder = "Database"
			t.Validate = utils.ValidateNotEmpty
		case 6:
			t.Placeholder = "Database driver"
			t.Validate = utils.ValidateDatabaseDriver
		}

		m.inputs[i] = t
	}

	return m
}

func (m addModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m addModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs) {
				return m, tea.Quit
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *addModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
		if err := m.inputs[i].Validate(m.inputs[i].Value()); err != nil {
			m.errors[i] = err.Error()
		} else {
			m.errors[i] = ""
		}
	}

	return tea.Batch(cmds...)
}

func (m addModel) View() string {
	var b strings.Builder

	for i := range m.inputs {
		inputView := m.inputs[i].View()
		errorView := ""
		if m.errors[i] != "" {
			errorView = errorStyle.Render(m.errors[i])
		}

		inputWidth := lipgloss.Width(inputView)

		padding := 120 - inputWidth - 100
		if padding < 0 {
			padding = 0
		}

		row := fmt.Sprintf("%s%s%s", inputView, strings.Repeat(" ", padding), errorView)
		b.WriteString(row)

		if i < len(m.inputs)-1 {
			b.WriteString("\n")
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(cursorModeHelpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(cursorModeHelpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new database configuration",

	Run: func(cmd *cobra.Command, args []string) {
		m := initialModel()

		addStart := tea.NewProgram(m)
		result, err := addStart.Run()
		if err != nil {
			utils.Log.Error("Could not start program", zap.Error(err))
			os.Exit(1)
		}

		m = result.(addModel)

		for _, input := range m.inputs {
			if err := input.Validate(input.Value()); err != nil {
				return
			}
		}

		configName := m.inputs[0].Value()
		host := m.inputs[1].Value()
		port := m.inputs[2].Value()
		user := m.inputs[3].Value()
		password := m.inputs[4].Value()
		database := m.inputs[5].Value()
		databaseDriver := m.inputs[6].Value()

		newConfig := config.DBConfig{
			ID:         uuid.New(),
			ConfigName: configName,
			Driver:     databaseDriver,
			Host:       host,
			Port:       port,
			User:       user,
			Password:   password,
			Database:   database,
		}

		config.Configs = append(config.Configs, newConfig)
		if err := utils.SaveConfigs(config.Configs, config.ConfigFile); err != nil {
			utils.Log.Error("Failed to save configuration", zap.Error(err))
		} else {
			fmt.Println("Configuration saved successfully.")
		}
	},
}

func init() {
	addCmd.Flags().BoolP("help", "h", false, "help for add")
	addCmd.Flags().MarkHidden("help")
}
