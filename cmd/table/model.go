/*
Copyright © 2024 Denys <https://github.com/AnyoneClown>
This is my license. There are many like it, but this one is mine.
My license is my best friend. It is my life. I must master it as I must
master my life.
*/
package table

import (
	"fmt"
	"strings"

	"github.com/AnyoneClown/anydb/utils"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	db          *sqlx.DB
	list        list.Model
	table       table.Model
	tableChosen bool
	chosenTable string
	limit       int
	width       int
	height      int
}

type Item struct {
	TableName string
	RowsCount int
}

func (i Item) Title() string       { return i.TableName }
func (i Item) Description() string { return fmt.Sprintf("Rows: %d", i.RowsCount) }
func (i Item) FilterValue() string { return i.TableName }

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.width = msg.Width-h
		m.height = msg.Height-v
		m.list.SetSize(m.width, m.height)

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if m.tableChosen {
				m.tableChosen = false
				m.list.SetSize(m.width, m.height)
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if !m.tableChosen {
				m.chosenTable = m.list.SelectedItem().(Item).TableName
				m.tableChosen = true
				var err error
				m.table, err = initializeTableData(m.db, m.chosenTable, m.limit)
				if err != nil {
					utils.Log.Error("Failed to initialize table data", zap.Error(err))
					return m, tea.Quit
				}
			}
		}
	}

	var cmd tea.Cmd
	if m.tableChosen {
		m.table, cmd = m.table.Update(msg)
	} else {
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	if m.tableChosen {
		return baseStyle.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n"
	}
	return docStyle.Render(m.list.View())
}

func initializeTableList(db *sqlx.DB) (list.Model, error) {
	tables, err := utils.GetTables(db)
	if err != nil {
		utils.Log.Error("Failed retrieve tables", zap.Error(err))
		return list.Model{}, err
	}

	items := make([]list.Item, len(tables))
	for i, table := range tables {
		items[i] = Item{TableName: table.TableName, RowsCount: table.RowsCount}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(lipgloss.Color("170")).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("170"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(lipgloss.Color("241"))

	resultList := list.New(items, delegate, 0, 0)
	resultList.Title = "Choose Database"
	resultList.SetShowStatusBar(true)
	resultList.SetFilteringEnabled(true)
	resultList.Styles.Title = lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1)

	return resultList, nil
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
	l, err := initializeTableList(db)
	if err != nil {
		return model{}, err
	}

	m := model{
		db:          db,
		list:        l,
		tableChosen: false,
		chosenTable: "",
		limit:       limit,
		width:       0,
		height:      0,
	}

	return m, nil
}
