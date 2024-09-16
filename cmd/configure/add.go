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
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type addModel struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
	errors []string
}

func validateNotEmpty(value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("field cannot be empty")
	}
	return nil
}

func validatePort(value string) error {
	port, err := strconv.Atoi(value)
	if err != nil || port <= 0 || port > 65535 {
		return fmt.Errorf("invalid port number")
	}
	return nil
}


func initialModel() addModel {
	m := addModel{
		inputs: make([]textinput.Model, 6),
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
			t.Validate = validateNotEmpty
		case 1:
			t.Placeholder = "Host"
			t.CharLimit = 64
			t.Validate = validateNotEmpty
		case 2:
			t.Placeholder = "Port"
			t.CharLimit = 64
			t.Validate = validatePort
		case 3:
			t.Placeholder = "User"
			t.CharLimit = 64
			t.Validate = validateNotEmpty
		case 4:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
			t.Validate = validateNotEmpty
		case 5:
			t.Placeholder = "Database"
			t.Validate = validateNotEmpty
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

		// Change cursor mode
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

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				return m, tea.Quit
			}

			// Cycle indexes
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
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *addModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m addModel) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new database configuration",

	Run: func(cmd *cobra.Command, args []string) {
		addStart := tea.NewProgram(initialModel())
		result, err := addStart.Run()
		if  err != nil {
			fmt.Printf("could not start program: %s\n", err)
			os.Exit(1)
		}
		
		m := result.(addModel)
		var validationErrors []string

		for _, input := range m.inputs {
			if err := input.Validate(input.Value()); err != nil {
				validationErrors = append(validationErrors, err.Error())
			}
		}

		if len(validationErrors) > 0 {
			fmt.Println("Validation errors:")
			for index, vErr := range validationErrors {
				fmt.Printf("%d: %s\n", index + 1, vErr)
			}
		}

		name := m.inputs[0].Value()
		host := m.inputs[1].Value()
		port := m.inputs[2].Value()
		user := m.inputs[3].Value()
		password := m.inputs[4].Value()
		database := m.inputs[5].Value()

		portInt, err := strconv.Atoi(port)
		if err != nil {
			return
		}

		config := DBConfig{
			Name:     name,
			Host:     host,
			Port:     portInt,
			User:     user,
			Password: password,
			Database: database,
		}

		configs = append(configs, config)
		if err := saveConfigs(configs, configFile); err != nil {
			fmt.Printf("Failed to save configuration: %v\n", err)
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

	addCmd.Flags().String("name", "", "Configuration name")
	addCmd.Flags().String("host", "localhost", "Database host")
	addCmd.Flags().Int("port", 5432, "Database port")
	addCmd.Flags().String("user", "", "Database user")
	addCmd.Flags().String("password", "", "Database password")
	addCmd.Flags().String("database", "", "Database name")

	// addCmd.MarkFlagRequired("name")
	// addCmd.MarkFlagRequired("host")
	// addCmd.MarkFlagRequired("port")
	// addCmd.MarkFlagRequired("user")
	// addCmd.MarkFlagRequired("password")
	// addCmd.MarkFlagRequired("database")

	addCmd.Flags().BoolP("help", "h", false, "help for add")
	addCmd.Flags().MarkHidden("help")
}
