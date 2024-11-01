package utils

import (
	"fmt"
	"os"

	"github.com/AnyoneClown/anydb/config"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

func LoadConfigs(file string) ([]config.DBConfig, error) {
	var configs []config.DBConfig

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

func SaveConfigs(configs []config.DBConfig, file string) error {
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
	data, err := os.ReadFile(config.DefaultConfigFile)
	if err != nil {
		Log.Error("Failed to read default configuration file", zap.Error(err))
		return err
	}

	err = yaml.Unmarshal(data, &config.DefaultConfigData)
	if err != nil {
		Log.Error("Failed to unmarshal default configuration data", zap.Error(err))
		return err
	}

	return nil
}

func GetConfigByID(id uuid.UUID) (*config.DBConfig, error) {
	configs, err := LoadConfigs(config.ConfigFile)
	if err != nil {
		return nil, err
	}

	for _, cfg := range configs {
		if cfg.ID == id {
			return &cfg, nil
		}
	}

	Log.Error("Configuration not found", zap.String("id", id.String()))
	return nil, fmt.Errorf("configuration with ID %s not found", id)
}
