package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env      string        `yaml:"env"`
	TokenTTL time.Duration `yaml:"token_ttl"`
	Postgres `yaml:"postgres"`
	GRPC     `yaml:"grpc"`
}

type Postgres struct {
	Host     string `yaml:"host" env-required:"true""`
	Port     int    `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	Name     string `yaml:"name" env-required:"true"`
}

type GRPC struct {
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-required:"true"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config file path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found")
	}

	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("SSO_CONFIG_PATH")
	}

	return res
}
