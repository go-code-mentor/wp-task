package app

import (
	"github.com/go-code-mentor/wp-task/internal/handlers"
	"github.com/go-code-mentor/wp-task/internal/service"
	"github.com/gofiber/fiber/v2"
)

func New(cfg Config) *App {
	return &App{
		cfg: cfg,
	}
}

type App struct {
	cfg    Config
	server *fiber.App
}

func (a *App) Build() error {

	a.server = fiber.New()

	appService := service.New(&service.FakeStorage{})
	tasksHandler := handlers.TasksHandler{Service: appService}

	api := a.server.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/tasks", tasksHandler.ListHandler)
	v1.Get("/tasks/:id", tasksHandler.ItemHandler)
	v1.Post("/tasks", tasksHandler.AddHandler)

	return nil
}

func (a *App) Run() error {
	return a.server.Listen(":3000")
}
