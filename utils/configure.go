/*
Copyright © 2024 Denys <https://github.com/AnyoneClown>
This is my license. There are many like it, but this one is mine.
My license is my best friend. It is my life. I must master it as I must
master my life.
*/
package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/AnyoneClown/anydb/config"
	"go.uber.org/zap"
)

func CreateFileAndDir() error {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".anydb")
	config.ConfigFile = filepath.Join(configDir, "anydb-config.yaml")
	config.DefaultConfigFile = filepath.Join(configDir, "anydb-default-config.yaml")

	// Check if the directory exists, if not, create it
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err = os.Mkdir(configDir, os.ModePerm)
		if err != nil {
			Log.Error("Failed to create configDir", zap.String("configDir", configDir), zap.Error(err))
			return err
		}
	}

	// Check if the config file exists, if not, create it
	if _, err := os.Stat(config.ConfigFile); os.IsNotExist(err) {
		file, err := os.Create(config.ConfigFile)
		if err != nil {
			Log.Error("Failed to create ConfigFile", zap.String("ConfigFile", config.ConfigFile), zap.Error(err))
			return err
		}
		file.Close()
	}

	// Check if the default config file exists, if not, create it
	if _, err := os.Stat(config.DefaultConfigFile); os.IsNotExist(err) {
		file, err := os.Create(config.DefaultConfigFile)
		if err != nil {
			Log.Error("Failed to create DefaultConfigFile", zap.String("DefaultConfigFile", config.DefaultConfigFile), zap.Error(err))
			return err
		}
		file.Close()
	}

	return nil
}

func ValidateNotEmpty(value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("field cannot be empty")
	}
	return nil
}

func ValidatePort(value string) error {
	port, err := strconv.Atoi(value)
	if err != nil || port <= 0 || port > 65535 {
		return fmt.Errorf("invalid port number")
	}
	return nil
}

func ValidateDatabaseDriver(value string) error {
	value = strings.ToLower(value)
	for _, driver := range config.SupportedDrivers {
		if value == driver {
			return nil
		}
	}
	errorMessage := fmt.Sprintf("Supported drivers: %s", strings.Join(config.SupportedDrivers, ", "))
	return errors.New(errorMessage)
}
