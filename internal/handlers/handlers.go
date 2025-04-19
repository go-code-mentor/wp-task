package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/go-code-mentor/wp-task/internal/entities"
)

type Service interface {
	Tasks(ctx context.Context) ([]entities.Task, error)
	TaskAdd(ctx context.Context, task entities.Task) error
	Task(ctx context.Context, id uint64) (entities.Task, error)
	TaskRemove(ctx context.Context, id uint64) error
}

type TasksHandler struct {
	Service Service
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

func (h *TasksHandler) ListHandler(c *fiber.Ctx) error {

	tasks, err := h.Service.Tasks(c.Context())
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(tasks)
}

func (h *TasksHandler) AddHandler(c *fiber.Ctx) error {
	//Read body and parse JSON to DTO
	var task entities.Task
	err := json.Unmarshal(c.Body(), &task)
	if err != nil {
		return fiber.ErrBadRequest
	}

	//Add task with service
	err = h.Service.TaskAdd(c.Context(), task)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return nil
}

func (h *TasksHandler) RemoveHandler(c *fiber.Ctx) error {
	//Fetching task id from url
	taskId, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}

	//Removing task with service
	err = h.Service.TaskRemove(c.Context(), taskId)
	if errors.Is(err, entities.ErrNoTask) {
		return fiber.ErrNotFound
	}
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return nil
}
