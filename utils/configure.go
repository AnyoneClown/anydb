/*
Copyright © 2024 Denys <https://github.com/AnyoneClown>
This is my license. There are many like it, but this one is mine.
My license is my best friend. It is my life. I must master it as I must
master my life.
*/
package utils

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
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
			Log.Error("Failed to create configDir", zap.String("configDir", configDir), zap.Error(err))
			return err
		}
	}

	// Check if the config file exists, if not, create it
	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
		file, err := os.Create(ConfigFile)
		if err != nil {
			Log.Error("Failed to create ConfigFile", zap.String("ConfigFile", ConfigFile), zap.Error(err))
			return err
		}
		file.Close()
	}

	// Check if the default config file exists, if not, create it
	if _, err := os.Stat(DefaultConfigFile); os.IsNotExist(err) {
		file, err := os.Create(DefaultConfigFile)
		if err != nil {
			Log.Error("Failed to create DefaultConfigFile", zap.String("DefaultConfigFile", DefaultConfigFile), zap.Error(err))
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
		Log.Error("Failed to read configuration file", zap.Error(err))
		return nil, err
	}

	err = yaml.Unmarshal(data, &configs)
	if err != nil {
		Log.Error("Failed to unmarshal configuration data", zap.Error(err))
		return nil, err
	}

	return configs, nil
}

func SaveConfigs(configs []DBConfig, file string) error {
	data, err := yaml.Marshal(configs)
	if err != nil {
		Log.Error("Failed to marshal configuration data", zap.Error(err))
		return err
	}

	err = os.WriteFile(file, data, 0644)
	if err != nil {
		Log.Error("Failed to write configuration file", zap.Error(err))
		return err
	}

	return nil
}

func LoadDefaultConfig() error {
	data, err := os.ReadFile(DefaultConfigFile)
	if err != nil {
		Log.Error("Failed to read default configuration file", zap.Error(err))
		return err
	}

	err = yaml.Unmarshal(data, &DefaultConfigData)
	if err != nil {
		Log.Error("Failed to unmarshal default configuration data", zap.Error(err))
		return err
	}

	return nil
}
