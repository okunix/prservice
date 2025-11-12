package config

import (
	"os"

	"github.com/goccy/go-yaml"
)

var (
	conf Config
)

type PostgresConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Addr     string `yaml:"addr"`
	SSLMode  string `yaml:"sslMode"`
	DB       string `yaml:"db"`
}

type Config struct {
	Addr           string         `yaml:"addr"`
	PostgresConfig PostgresConfig `yaml:"postgres"`
}

func newDefaultConfig() Config {
	return Config{
		Addr: "0.0.0.0:80",
		PostgresConfig: PostgresConfig{
			User:     "postgres",
			Password: "postgres",
			Addr:     "localhost:5432",
			SSLMode:  "disable",
			DB:       "prservice_db",
		},
	}
}

func Read(path string) (Config, error) {
	defaultConfig := newDefaultConfig()
	file, err := os.Open(path)
	if err != nil {
		return defaultConfig, err
	}
	err = yaml.NewDecoder(file).Decode(&defaultConfig)
	conf = defaultConfig
	return defaultConfig, err
}

func GetConfig() Config {
	return conf
}
