/*
Copyright © 2024 Denys <https://github.com/AnyoneClown>
This is my license. There are many like it, but this one is mine.
My license is my best friend. It is my life. I must master it as I must
master my life.
*/
package configure

import (
	"os"

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
