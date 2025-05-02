package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server   *Server
	Postgres *Postgres
	Email    *Email
	Minio    *Minio
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

type Email struct {
	SmtpHost      string `env:"SMTP_HOST"`
	SmtpPort      int    `env:"SMTP_PORT"`
	EmailUser     string `env:"EMAIL_USER"`
	EmailPassword string `env:"EMAIL_PASSWORD"`
	AdminEmail    string `env:"ADMIN_EMAIl"`
}

type Minio struct {
	AccessKeyMinio string `env:"ACCESS_KEY_MINIO"`
	SecretMinio    string `env:"SECRET_MINIO"`
	MinioPort      int    `env:"MINIO_PORT"`
	MinioHost      string `env:"MINIO_HOST"`
	MinioBucket    string `env:"MINIO_BUCKET"`
}

// MustLoad - загрузка данных из .env в конфиг.
func MustLoad() Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is not specified")
	}
	serverConfig := Server{}
	postgresConfig := Postgres{}
	emailConfig := Email{}
	minioConfig := Minio{}
	cfg := Config{}
	if err := cleanenv.ReadConfig(configPath, &serverConfig); err != nil {
		panic(fmt.Sprintf("cannot read config from file; err=%s", err.Error()))
	}
	if err := cleanenv.ReadConfig(configPath, &postgresConfig); err != nil {
		panic(fmt.Sprintf("cannot read config from file; err=%s", err.Error()))
	}
	if err := cleanenv.ReadConfig(configPath, &emailConfig); err != nil {
		panic(fmt.Sprintf("cannot read config from file; err=%s", err.Error()))
	}
	if err := cleanenv.ReadConfig(configPath, &minioConfig); err != nil {
		panic(fmt.Sprintf("cannot read config from file; err=%s", err.Error()))
	}
	cfg.Server = &serverConfig
	cfg.Postgres = &postgresConfig
	cfg.Email = &emailConfig
	cfg.Minio = &minioConfig
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
