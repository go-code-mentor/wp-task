package app

import (
	"errors"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

func ParseConfig() (Config, error) {
	cfg := Config{}

	if err := cfg.parseDb(); err != nil {
		return cfg, err
	}

	return cfg, nil
}

type Config struct {
	pg_uri string
}

type ConfigDatabase struct {
	Port     string `yaml:"port" env:"POSTGRES_PORT" env-default:"5432"`
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-default:"localhost"`
	Name     string `yaml:"name" env:"POSTGRES_DB" env-default:"work_planner"`
	User     string `yaml:"user" env:"POSTGRES_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD"`
}

func (c *Config) parseDb() error {

	cfgFileName := ".env"

	var cfg ConfigDatabase

	if _, err := os.Stat(cfgFileName); err == nil {
		if err := cleanenv.ReadConfig(cfgFileName, &cfg); err != nil {
			return err
		}
	} else if errors.Is(err, os.ErrNotExist) {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			return err
		}
	}

	c.pg_uri = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	return nil
}
