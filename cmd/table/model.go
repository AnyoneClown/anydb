package table

import (
	"fmt"
	"strings"

	"github.com/AnyoneClown/anydb/utils"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jmoiron/sqlx"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	db          *sqlx.DB
	table       table.Model
	tableChosen bool
	chosenTable string
	limit       int
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.tableChosen {
				m.tableChosen = false
				var err error
				m.table, err = initializeTableList(m.db)
				if err != nil {
					return m, tea.Quit
				}
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if !m.tableChosen {
				m.chosenTable = m.table.SelectedRow()[0]
				m.tableChosen = true
				var err error
				m.table, err = initializeTableData(m.db, m.chosenTable, m.limit)
				if err != nil {
					return m, tea.Quit
				}
			}
		}
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n"
}

func initializeTableList(db *sqlx.DB) (table.Model, error) {
	tables, err := utils.GetTables(db)
	if err != nil {
		return table.Model{}, err
	}

	rows := make([]table.Row, len(tables))
	for i, tableName := range tables {
		rows[i] = table.Row{tableName}
	}

	columns := []table.Column{
		{Title: "Table Name", Width: 30},
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

	return t, nil
}

func initializeTableData(db *sqlx.DB, tableName string, limit int) (table.Model, error) {
	columns, err := utils.GetTableColumns(db, tableName)
	if err != nil {
		return table.Model{}, err
	}

	records, err := utils.GetLastRecords(db, tableName, limit)
	if err != nil {
		return table.Model{}, err
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

	return t, nil
}

func NewModel(db *sqlx.DB, limit int) (model, error) {
	t, err := initializeTableList(db)
	if err != nil {
		return model{}, err
	}

	return model{
		db:          db,
		table:       t,
		tableChosen: false,
		chosenTable: "",
		limit:       limit,
	}, nil
}

