package utils

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func CreateFileAndDir() error {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".anydb")
	ConfigFile = filepath.Join(configDir, "anydb-config.yaml")
	DefaultConfigFile = filepath.Join(configDir, "anydb-default-config.yaml")

	// Check if the directory exists, if not, create it
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err = os.Mkdir(configDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// Check if the config file exists, if not, create it
	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
		file, err := os.Create(ConfigFile)
		if err != nil {
			return err
		}
		file.Close()
	}

	// Check if the default config file exists, if not, create it
	if _, err := os.Stat(DefaultConfigFile); os.IsNotExist(err) {
		file, err := os.Create(DefaultConfigFile)
		if err != nil {
			return err
		}
		file.Close()
	}

	return nil
}

func LoadConfigs(file string) ([]DBConfig, error) {
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

func SaveConfigs(configs []DBConfig, file string) error {
	data, err := yaml.Marshal(configs)
	if err != nil {
		return err
	}

	return os.WriteFile(file, data, 0644)
}

func LoadDefaultConfig() error {
	data, err := os.ReadFile(DefaultConfigFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &DefaultConfigData)
	if err != nil {
		return err
	}
	return nil
}
