package app

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

func ParseConfig() (Config, error) {
	cfg := Config{}

	if err := cfg.parseDb(); err != nil {
		return cfg, err
	}

	if err := cfg.parseTg(); err != nil {
		return cfg, err
	}

	return cfg, nil
}

type Config struct {
	pg_uri string
	tg_uri string
}

func (c *Config) ConnString() string {
	return c.pg_uri
}

type ConfigDatabase struct {
	Port     string `yaml:"port" env:"POSTGRES_PORT" env-default:"5432"`
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-default:"localhost"`
	Name     string `yaml:"name" env:"POSTGRES_DB" env-default:"work_planner"`
	User     string `yaml:"user" env:"POSTGRES_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD"`
}

func (c *Config) parseDb() error {

	var cfg ConfigDatabase
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return err
	}

	c.pg_uri = fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable&search_path=public",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	return nil
}

type ConfigTgService struct {
	Host string `yaml:"tg_service_host" env:"TG_SERVICE_HOST" env-default:"localhost"`
	Port string `yaml:"tg_service_port" env:"TG_SERVICE_PORT" env-default:"8000"`
}

func (c *Config) parseTg() error {

	var cfg ConfigTgService
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return err
	}

	c.tg_uri = fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	return nil
}
