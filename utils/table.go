/*
Copyright © 2024 Denys <https://github.com/AnyoneClown>
This is my license. There are many like it, but this one is mine.
My license is my best friend. It is my life. I must master it as I must
master my life.
*/
package utils

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/jmoiron/sqlx"
)

func GetDBString() (string, error) {
	err := LoadDefaultConfig()
	if err != nil {
		return "", err
	}

	var dsn string
	switch DefaultConfigData.Driver {
	case "postgres", "cockroachdb":
		dsn = fmt.Sprintf(
			"postgresql://%s:%s@%s:%d/%s?sslmode=verify-full",
			DefaultConfigData.User,
			DefaultConfigData.Password,
			DefaultConfigData.Host,
			DefaultConfigData.Port,
			DefaultConfigData.Database,
		)
	default:
		return "", fmt.Errorf("unsupported database driver: %s", DefaultConfigData.Driver)
	}
	return dsn, nil
}

func GetLastRecords(db *sqlx.DB, tableName string, limit int) ([]map[string]interface{}, error) {
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

func GetTableColumns(db *sqlx.DB, tableName string) ([]table.Column, error) {
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

func GetTables(db *sqlx.DB) ([]string, error) {
	query := "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'"
	var tables []string
	err := db.Select(&tables, query)
	return tables, err
}
