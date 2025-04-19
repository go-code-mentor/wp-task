package app

import (
	"fmt"
	"os"
)

func ParseConfig() (Config, error) {
	cfg := Config{}

	cfg.parseDb()

	return cfg, nil
}

type Config struct {
	pg_uri string
}

func (c *Config) parseDb() {
	host := os.Getenv("POSTGRES_HOST")
	db_name := os.Getenv("POSTGRES_DB")
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5432"
	}
	c.pg_uri = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, db_name)
}
