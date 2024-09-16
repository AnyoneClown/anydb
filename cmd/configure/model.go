package configure

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	dbConfig DBConfig
}

func (i item) Title() string { return i.dbConfig.Name }
func (i item) Description() string {
	return fmt.Sprintf("%s@%s:%s", i.dbConfig.Name, i.dbConfig.Host, i.dbConfig.Database)
}
func (i item) FilterValue() string { return i.dbConfig.Name }

type model struct {
	list      list.Model
	choice    *DBConfig
	quitting  bool
	noConfigs bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = &i.dbConfig
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != nil {
		return fmt.Sprintf("You chose: %s\n", m.choice.Name)
	}
	if m.quitting && !m.noConfigs {
		return "No configuration selected. Goodbye!\n"
	}
	if m.noConfigs {
		return docStyle.Render(`
	No configurations available.

	Press 'enter' or 'q' to exit

	You can also use 'anydb configure add' from the command line to add a new configuration.`)
	}
	return docStyle.Render(m.list.View())
}

func newModel() model {
	items := make([]list.Item, len(configs)) // Create slice with configuration items
	for i, config := range configs {
		items[i] = item{dbConfig: config}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(lipgloss.Color("170")).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("170"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(lipgloss.Color("241"))

	m := model{
		list:      list.New(items, delegate, 0, 0),
		noConfigs: len(configs) == 0,
	}
	m.list.Title = "Database Configurations"
	m.list.SetShowStatusBar(true)
	m.list.SetFilteringEnabled(true)
	m.list.Styles.Title = lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1)

	return m
}
