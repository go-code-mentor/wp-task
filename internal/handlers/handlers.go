package handlers

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/go-code-mentor/wp-task/internal/entities"
)

type Service interface {
	Tasks(ctx context.Context) ([]entities.Task, error)
	Task(ctx context.Context, id uint64) (entities.Task, error)
}

type TasksHandler struct {
	Service Service
}

func (h *TasksHandler) ListHandler(c *fiber.Ctx) error {

	tasks, err := h.Service.Tasks(c.Context())
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(tasks)
}

func (h *TasksHandler) ItemHandler(c *fiber.Ctx) error {

	taskId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.ErrNotFound
	}

	task, err := h.Service.Task(c.Context(), uint64(taskId))
	if err != nil {
		return fiber.ErrNotFound
	}

	return c.JSON(task)
}
