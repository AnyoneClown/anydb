package last

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jmoiron/sqlx"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
}

func getLastFiveRecords(db *sqlx.DB, tableName string, limit int) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id DESC LIMIT %d", tableName, limit)
	rows, err := db.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		result := make(map[string]interface{})
		err := rows.MapScan(result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n"
}

func initializeModel(db *sqlx.DB, tableName string, limit int) (model, error) {
	columns := []table.Column{
		{Title: "ID", Width: 5},
		{Title: "Email", Width: 25},
		{Title: "Picture", Width: 10},
		{Title: "Updated At", Width: 20},
		{Title: "Is Admin", Width: 10},
		{Title: "Deleted At", Width: 20},
		{Title: "Password", Width: 60},
		{Title: "Google ID", Width: 10},
		{Title: "Provider", Width: 10},
		{Title: "Verified Email", Width: 15},
		{Title: "Created At", Width: 20},
	}

	records, err := getLastFiveRecords(db, tableName, limit)
	if err != nil {
		return model{}, err
	}

	var rows []table.Row
	for _, record := range records {
		row := table.Row{
			fmt.Sprintf("%v", record["id"]),
			fmt.Sprintf("%v", record["email"]),
			fmt.Sprintf("%v", record["picture"]),
			fmt.Sprintf("%v", record["updated_at"]),
			fmt.Sprintf("%v", record["is_admin"]),
			fmt.Sprintf("%v", record["deleted_at"]),
			fmt.Sprintf("%v", record["password"]),
			fmt.Sprintf("%v", record["google_id"]),
			fmt.Sprintf("%v", record["provider"]),
			fmt.Sprintf("%v", record["verified_email"]),
			fmt.Sprintf("%v", record["created_at"]),
		}
		rows = append(rows, row)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return model{t}, nil
}
