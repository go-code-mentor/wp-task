package app

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/go-code-mentor/wp-task/internal/handlers"
	"github.com/go-code-mentor/wp-task/internal/middleware/simpletoken"
	"github.com/go-code-mentor/wp-task/internal/service"
	"github.com/go-code-mentor/wp-task/internal/service/tgclient"
	userservice "github.com/go-code-mentor/wp-task/internal/service/users"
	"github.com/go-code-mentor/wp-task/internal/storage"
	tgapi "github.com/go-code-mentor/wp-tg-bot/api"
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
	tgConn *grpc.ClientConn
}

func (a *App) Build() error {

	if err := a.connectDb(); err != nil {
		return fmt.Errorf("failed to get db connection: %w", err)
	}

	if err := a.connectTg(); err != nil {
		return fmt.Errorf("failed to get tg bot connection: %w", err)
	}

	a.server = fiber.New()

	appStorage := storage.New(a.conn)

	userService := userservice.New(appStorage)
	authMiddleware := simpletoken.AuthMiddleware{Service: userService}
	a.server.Use(authMiddleware.Auth)

	tgBot := tgapi.NewTgBotClient(a.tgConn)
	tgService := tgclient.New(tgBot)
	appService := service.New(appStorage, tgService)
	tasksHandler := handlers.TasksHandler{Service: appService}

	api := a.server.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/tasks", tasksHandler.ListHandler)
	v1.Get("/tasks/:id", tasksHandler.ItemHandler)
	v1.Put("/tasks/:id", tasksHandler.UpdateHandler)
	v1.Post("/tasks", tasksHandler.AddHandler)
	v1.Delete("/tasks/:id", tasksHandler.RemoveHandler)

	return nil
}

func (a *App) Run() error {
	defer func() {
		err := a.conn.Close(context.Background())
		if err != nil {
			log.Errorf("failed to close db connection: %s", err)
		}

		err = a.tgConn.Close()
		if err != nil {
			log.Errorf("failed to close tg bot connection: %s", err)
		}
	}()
	return a.server.Listen(":3000")
}

func (a *App) connectDb() error {

	conn, err := pgx.Connect(context.Background(), a.cfg.pg_uri)
	if err != nil {
		return fmt.Errorf("could not connect db: %w", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		return fmt.Errorf("could not ping db: %w", err)
	}

	a.conn = conn

	return nil
}

func (a *App) connectTg() error {
	conn, err := grpc.NewClient(a.cfg.tg_uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("could not connect tg bot: %w", err)
	}

	a.tgConn = conn

	return nil
}
