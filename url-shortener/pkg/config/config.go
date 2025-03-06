package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	HTTP
	DB
}

type HTTP struct {
	Port         string        `env:"HTTP_PORT" env-default:":8080"`
	IdleTimeout  time.Duration `env:"HTTP_IDLE_TIMEOUT" env-default:"30s"`
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"15s"`
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"15s"`
}

type DB struct {
	User string `env:"PG_USER" env-required:"true"`
	Pass string `env:"PG_PASS" env-required:"true"`
	Host string `env:"PG_HOST" env-required:"true"`
	Port int    `env:"PG_PORT" env-required:"true"`
	Name string `env:"PG_NAME" env-required:"true"`
}

func Load() (*Config, error) {
	var conf Config

	err := cleanenv.ReadEnv(&conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
