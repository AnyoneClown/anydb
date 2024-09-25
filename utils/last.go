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
