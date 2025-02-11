package main

import (
	"os"

	//nolint:depguard
	"github.com/BurntSushi/toml"
)

type Config struct {
	Logger   LoggerConf   `toml:"logger"`
	Storage  StorageConf  `toml:"storage"`
	Database DatabaseConf `toml:"database"`
	Server   ServerConf   `toml:"server"`
}

type LoggerConf struct {
	LogFile  string `toml:"logFile"`
	LogLevel string `toml:"logLevel"`
}

type StorageConf struct {
	StorageType string `toml:"storageType"`
}

type DatabaseConf struct {
	DBConnection string `toml:"dbConnection"`
	DBHost       string `toml:"dbHost"`
	DBPort       int    `toml:"dbPort"`
	DBDatabase   string `toml:"dbDatabase"`
	DBUsername   string `toml:"dbUsername"`
	DBPassword   string `toml:"dbPassword"`
}

type ServerConf struct {
	HTTPHost string `toml:"httpHost"`
	HTTPPort int    `toml:"httpPort"`
}

func NewConfig(configPath string) (Config, error) {
	var config Config

	pwd, err := os.Getwd()
	if err != nil {
		return Config{}, err
	}

	_, err = toml.DecodeFile(pwd+configPath, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
