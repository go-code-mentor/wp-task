package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/go-code-mentor/wp-task/internal/app"
)

func main() {

	cfg, err := app.ParseConfig()
	if err != nil {
		log.Fatalf("failed to pasre config: %s", err)
	}

	m, err := migrate.New("file://internal/db/migrations", cfg.ConnString())
	if err != nil {
		log.Fatal(err)
	}

	args := os.Args

	if len(args) < 2 {
		log.Fatal("Please provide the migration direction (up or down) as the second argument.")
	}

	if args[1] != "up" && args[1] != "down" {
		log.Fatalf("Unexpected the migration direction %q. Chioce up or down.", args[1])
	}

	if args[1] == "up" {
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
	}

	if args[1] == "down" {
		if err := m.Down(); err != nil {
			log.Fatal(err)
		}
	}
}
