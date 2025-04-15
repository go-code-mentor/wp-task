package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/go-code-mentor/wp-task/internal/app"
	"github.com/go-code-mentor/wp-task/internal/handlers"
	"github.com/go-code-mentor/wp-task/internal/service"
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
	}

	appService := service.New(&service.FakeStorage{})
	tasksHandler := handlers.TasksHandler{Service: appService}

	webApp := fiber.New()

	api := webApp.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/tasks", tasksHandler.ListHandler)
	v1.Get("/tasks/:id", tasksHandler.ItemHandler)

	log.Fatal(webApp.Listen(":3000"))
}
