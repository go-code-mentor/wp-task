package app

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"

	"github.com/go-code-mentor/wp-task/internal/handlers"
	"github.com/go-code-mentor/wp-task/internal/service"
)

func New(cfg Config) *App {
	return &App{
		cfg: cfg,
	}
}

type App struct {
	cfg    Config
	server *fiber.App
	conn   *pgx.Conn
}

func (a *App) Build() error {

	if err := a.connectDb(); err != nil {
		return err
	}

	a.server = fiber.New()

	appService := service.New(&service.FakeStorage{})
	tasksHandler := handlers.TasksHandler{Service: appService}

	api := a.server.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/tasks", tasksHandler.ListHandler)
	v1.Get("/tasks/:id", tasksHandler.ItemHandler)

	return nil
}

func (a *App) Run() error {
	defer a.conn.Close(context.Background())
	return a.server.Listen(":3000")
}

func (a *App) connectDb() error {

	conn, err := pgx.Connect(context.Background(), a.cfg.pg_uri)
	if err != nil {
		return fmt.Errorf("could not connect db: %w", err)
	}

	a.conn = conn

	return nil
}
