package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/go-code-mentor/wp-task/internal/entities"
	"github.com/go-code-mentor/wp-task/internal/middleware/simpletoken"
)

type Service interface {
	Tasks(ctx context.Context, login string) ([]entities.Task, error)
	TaskAdd(ctx context.Context, task entities.Task, login string) error
	Task(ctx context.Context, id uint64, login string) (entities.Task, error)
	TaskRemove(ctx context.Context, id uint64, login string) error
	TaskUpdate(ctx context.Context, task entities.Task, login string) error
}

type TasksHandler struct {
	Service Service
}

type TaskJSON struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *TasksHandler) ItemHandler(c *fiber.Ctx) error {

	login := c.Locals(simpletoken.UserLoginKey).(string)

	taskId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.ErrNotFound
	}

	task, err := h.Service.Task(c.Context(), uint64(taskId), login)
	if err != nil {
		return fiber.ErrNotFound
	}

	return c.JSON(task)
}

func (h *TasksHandler) ListHandler(c *fiber.Ctx) error {

	login := c.Locals(simpletoken.UserLoginKey).(string)

	tasks, err := h.Service.Tasks(c.Context(), login)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(tasks)
}

func (h *TasksHandler) AddHandler(c *fiber.Ctx) error {

	login := c.Locals(simpletoken.UserLoginKey).(string)

	//Read body and parse JSON to DTO
	var task entities.Task
	err := json.Unmarshal(c.Body(), &task)
	if err != nil {
		return fiber.ErrBadRequest
	}

	//Add task with service
	err = h.Service.TaskAdd(c.Context(), task, login)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *TasksHandler) RemoveHandler(c *fiber.Ctx) error {

	login := c.Locals(simpletoken.UserLoginKey).(string)

	//Fetching task id from url
	taskId, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}

	//Removing task with service
	err = h.Service.TaskRemove(c.Context(), taskId, login)
	if errors.Is(err, entities.ErrNoTask) {
		return fiber.ErrNotFound
	}
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return nil
}

func (h *TasksHandler) UpdateHandler(c *fiber.Ctx) error {

	login := c.Locals(simpletoken.UserLoginKey).(string)

	taskId, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}

	//Read body and parse JSON to DTO
	var taskDTO TaskJSON
	err = json.Unmarshal(c.Body(), &taskDTO)
	if err != nil {
		return fiber.ErrBadRequest
	}

	//Convert to task entity
	task := entities.Task{
		ID:          uint64(taskId),
		Name:        taskDTO.Name,
		Description: taskDTO.Description,
	}

	//Update task in service
	err = h.Service.TaskUpdate(c.Context(), task, login)
	if errors.Is(err, entities.ErrNoTask) {
		return fiber.ErrNotFound
	}
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusNoContent)
}
