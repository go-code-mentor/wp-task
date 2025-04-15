package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/go-code-mentor/wp-task/internal/entities"
)

type Service interface {
	Tasks() ([]entities.Task, error)
	Task(id string) (entities.Task, error)
}

type TasksHandler struct {
	Service Service
}

func (h *TasksHandler) ListHandler(c *fiber.Ctx) error {

	tasks, err := h.Service.Tasks()
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(tasks)
}

func (h *TasksHandler) ItemHandler(c *fiber.Ctx) error {

	task, err := h.Service.Task(c.Params("id"))
	if err != nil {
		return fiber.ErrNotFound
	}

	return c.JSON(task)
}
