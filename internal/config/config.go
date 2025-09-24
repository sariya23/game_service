package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	TestEnvType = "test"
)

type Config struct {
	Server   *Server
	Postgres *Postgres
	Minio    *Minio
	Env      *Env
}

type Env struct {
	EnvType string `env:"ENV_TYPE"`
}

type Server struct {
	GrpcServerPort       int `env:"GRPC_SERVER_PORT"`
	HttpServerPort       int `env:"HTTP_SERVER_PORT"`
	ServerTimeoutSeconds int `env:"SERVER_TIMEOUT_SECONDS"`
}

type Postgres struct {
	PostgresURL      string `env:"POSTGRES_URL"`
	PostgresDBName   string `env:"POSTGRES_DB"`
	PostgresUsername string `env:"POSTGRES_USERNAME"`
	PostgresPassword string `env:"POSTRGRES_PASSWORD"`
}

type Minio struct {
	MinioUser     string `env:"MINIO_USER"`
	MinioPassword string `env:"MINIO_PASSWORD"`
	MinioPort     int    `env:"MINIO_PORT"`
	MinioHost     string `env:"MINIO_HOST"`
	MinioBucket   string `env:"MINIO_BUCKET"`
}

// MustLoad - загрузка данных из .env в конфиг.
func MustLoad() Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is not specified")
	}
	serverConfig := Server{}
	postgresConfig := Postgres{}
	minioConfig := Minio{}
	envConfig := Env{}
	cfg := Config{}
	if err := cleanenv.ReadConfig(configPath, &serverConfig); err != nil {
		panic(fmt.Sprintf("cannot read config from file; err=%s", err.Error()))
	}
	if err := cleanenv.ReadConfig(configPath, &postgresConfig); err != nil {
		panic(fmt.Sprintf("cannot read config from file; err=%s", err.Error()))
	}
	if err := cleanenv.ReadConfig(configPath, &minioConfig); err != nil {
		panic(fmt.Sprintf("cannot read config from file; err=%s", err.Error()))
	}
	if err := cleanenv.ReadConfig(configPath, &envConfig); err != nil {
		panic(fmt.Sprintf("cannot read config from file; err=%s", err.Error()))
	}
	cfg.Server = &serverConfig
	cfg.Postgres = &postgresConfig
	cfg.Minio = &minioConfig
	cfg.Env = &envConfig
	return cfg
}

// MustLoadByPath - загрузка конфига по пути.
func MustLoadByPath(configPath string) Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exists: " + configPath)
	}

	serverConfig := Server{}
	postgresConfig := Postgres{}
	minioConfig := Minio{}
	envConfig := Env{}
	cfg := Config{}
	if err := cleanenv.ReadConfig(configPath, &serverConfig); err != nil {
		panic(fmt.Sprintf("cannot read config from file; err=%s", err.Error()))
	}
	if err := cleanenv.ReadConfig(configPath, &postgresConfig); err != nil {
		panic(fmt.Sprintf("cannot read config from file; err=%s", err.Error()))
	}
	if err := cleanenv.ReadConfig(configPath, &minioConfig); err != nil {
		panic(fmt.Sprintf("cannot read config from file; err=%s", err.Error()))
	}
	if err := cleanenv.ReadConfig(configPath, &envConfig); err != nil {
		panic(fmt.Sprintf("cannot read config from file; err=%s", err.Error()))
	}
	cfg.Server = &serverConfig
	cfg.Postgres = &postgresConfig
	cfg.Minio = &minioConfig
	cfg.Env = &envConfig
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
