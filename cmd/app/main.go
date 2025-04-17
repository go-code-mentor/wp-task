package main

import (
	"log"

	"github.com/go-code-mentor/wp-task/internal/app"
)

func main() {
	cfg, err := app.ParseConfig()
	if err != nil {
		log.Fatalf("failed to pasre config: %s", err)
	}

	a := app.New(cfg)

	if err := a.Build(); err != nil {
		log.Fatalf("failed to build app: %s", err)
	}

	if err := a.Run(); err != nil {
		log.Fatalf("failed to run app: %s", err)
	} else {
		log.Println("app successfully stopped")
	}

}
