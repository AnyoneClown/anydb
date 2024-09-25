package last

import (
	"fmt"
	"strings"

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

func getLastRecords(db *sqlx.DB, tableName string, limit int) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id DESC LIMIT %d", tableName, limit)
	rows, err := db.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]map[string]interface{}, 0, limit)
	for rows.Next() {
		result := make(map[string]interface{})
		if err := rows.MapScan(result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, rows.Err()
}

func getTableColumns(db *sqlx.DB, tableName string) ([]table.Column, error) {
	query := "SELECT column_name FROM information_schema.columns WHERE table_name = $1"
	var columnNames []string
	if err := db.Select(&columnNames, query, tableName); err != nil {
		return nil, err
	}

	columns := make([]table.Column, len(columnNames))
	for i, name := range columnNames {
		columns[i] = table.Column{Title: name, Width: 20}
	}

	return columns, nil
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.table.Focus()
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Printf("Let's go to %s!", m.table.SelectedRow()[1])
		}
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n"
}

func initializeModel(db *sqlx.DB, tableName string, limit int) (model, error) {
	columns, err := getTableColumns(db, tableName)
	if err != nil {
		return model{}, err
	}

	records, err := getLastRecords(db, tableName, limit)
	if err != nil {
		return model{}, err
	}

	rows := make([]table.Row, len(records))
	for i, record := range records {
		row := make(table.Row, len(columns))
		for j, column := range columns {
			row[j] = fmt.Sprintf("%v", record[strings.ToLower(column.Title)])
		}
		rows[i] = row
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
