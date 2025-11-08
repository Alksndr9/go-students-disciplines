package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env"         env-required:"true"`

	PG         `yaml:"pg"          env-required:"true"`
	HTTPServer `yaml:"http_server"          env-required:"true"`
}

type PG struct {
	User     string `yaml:"user"     env:"PG_USER"     env-required:"true"`
	Password string `yaml:"password" env:"PG_PASSWORD" env-required:"true"`
	Host     string `yaml:"host"     env:"PG_HOST"     env-required:"true"`
	Port     string `yaml:"port"     env:"PG_PORT"     env-required:"true"`
	Database string `yaml:"database" env:"PG_DATABASE" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `yaml:"address"      env-required:"true"`
	Timeout     time.Duration `yaml:"timeout"      env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("cofig file is not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("connot read config: %s", err)
	}

	return &cfg
}
