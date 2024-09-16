/*
Copyright © 2024 Denys <https://github.com/AnyoneClown>
This is my license. There are many like it, but this one is mine.
My license is my best friend. It is my life. I must master it as I must
master my life.
*/
package configure

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type DBConfig struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

var configs []DBConfig
var configFile string
var defaultConfigFile string

func createFileAndDir() error {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".anydb")
	configFile = filepath.Join(configDir, "anydb-config.yaml")
	defaultConfigFile = filepath.Join(configDir, "anydb-default-config.yaml")

	// Check if the directory exists, if not, create it
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err = os.Mkdir(configDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// Check if the config file exists, if not, create it
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		file, err := os.Create(configFile)
		if err != nil {
			return err
		}
		file.Close()
	}

	// Check if the default config file exists, if not, create it
	if _, err := os.Stat(defaultConfigFile); os.IsNotExist(err) {
		file, err := os.Create(defaultConfigFile)
		if err != nil {
			return err
		}
		file.Close()
	}

	return nil
}

func loadConfigs(file string) ([]DBConfig, error) {
	var configs []DBConfig

	data, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return configs, nil
		}
		return nil, err
	}

	err = yaml.Unmarshal(data, &configs)
	if err != nil {
		return nil, err
	}

	return configs, nil
}

func saveConfigs(configs []DBConfig, file string) error {
	data, err := yaml.Marshal(configs)
	if err != nil {
		return err
	}

	return os.WriteFile(file, data, 0644)
}

func containsErrors(errors []string) bool {
	for _, err := range errors {
		if err != "" {
			return true
		}
	}
	return false
}