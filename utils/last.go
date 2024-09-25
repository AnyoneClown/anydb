/*
Copyright © 2024 Denys <https://github.com/AnyoneClown>
This is my license. There are many like it, but this one is mine.
My license is my best friend. It is my life. I must master it as I must
master my life.
*/
package utils

import "fmt"

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
