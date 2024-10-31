/*
Copyright © 2024 Denys <https://github.com/AnyoneClown>
This is my license. There are many like it, but this one is mine.
My license is my best friend. It is my life. I must master it as I must
master my life.
*/
package config

import "github.com/google/uuid"

type DBConfig struct {
	ID         uuid.UUID `yaml:"id"`
	ConfigName string    `yaml:"configName"`
	Driver     string    `yaml:"driver"`
	Host       string    `yaml:"host"`
	Port       string    `yaml:"port"`
	User       string    `yaml:"user"`
	Password   string    `yaml:"password"`
	Database   string    `yaml:"database"`
}

type ConfigInput struct {
	ConfigName string `json:"configName"`
	Driver     string `json:"driver"`
	Host       string `json:"host"`
	Port       string `json:"port"`
	User       string `json:"user"`
	Password   string `json:"password"`
	Database   string `json:"database"`
}

var Configs []DBConfig
var ConfigFile string
var DefaultConfigFile string
var DefaultConfigData DBConfig

var SupportedDrivers = []string{
	"postgres",
	"cockroachdb",
}
