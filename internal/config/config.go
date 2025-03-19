package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config - конфиг всего приложения.
type Config struct {
	GrpcServerPort       int    `env:"GRPC_SERVER_PORT"`
	HttpServerPort       int    `env:"HTTP_SERVER_PORT"`
	ServerTimeoutSeconds int    `env:"SERVER_TIMEOUT_SECONDS"`
	PostgresURL          string `env:"POSTGRES_URL"`
	PostgresDBName       string `env:"POSTGRES_DB"`
	PostgresUsername     string `env:"POSTGRES_USERNAME"`
	PostgresPassword     string `env:"POSTRGRES_PASSWORD"`
}

// MustLoad - загрузка данных из .env в конфиг.
func MustLoad() Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is not specified")
	}
	cfg := Config{}
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(fmt.Sprintf("cannot read config from file; err=%s", err.Error()))
	}
	return cfg
}

// MustLoadByPath - загрузка конфига по пути.
func MustLoadByPath(path string) Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exists: " + path)
	}

	cfg := Config{}
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}
	return cfg
}

// fetchConfigPath - парсит пусть до файла с конфигом.
// Приоритет: значение из флага при запуске > дефолтное значение.
func fetchConfigPath() string {
	var configPath string

	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.Parse()
	return configPath
}
