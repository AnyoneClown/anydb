/*
Copyright © 2024 Denys <https://github.com/AnyoneClown>
This is my license. There are many like it, but this one is mine.
My license is my best friend. It is my life. I must master it as I must
master my life.
*/
package config

type DBConfig struct {
	ConfigName string `yaml:"configName"`
	Driver     string `yaml:"driver"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Database   string `yaml:"database"`
}

var Configs []DBConfig
var ConfigFile string
var DefaultConfigFile string
var DefaultConfigData DBConfig

var SupportedDrivers = []string{
	"postgres",
	"cockroachdb",
}
