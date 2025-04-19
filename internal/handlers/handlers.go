package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/go-code-mentor/wp-task/internal/entities"
	"io"
	"net/http"
	"strings"
)

type Service interface {
	Tasks(ctx context.Context) ([]entities.Task, error)
	TaskAdd(ctx context.Context, task entities.Task) error
	Task(ctx context.Context, id uint64) (entities.Task, error)
	TaskRemove(ctx context.Context, id uint64) error
	TaskUpdate(ctx context.Context, task entities.Task) error
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

func (h *TasksHandler) AddHandler(w http.ResponseWriter, r *http.Request) {
	// Check HTTP Method
	if r.Method != http.MethodPost {
		ErrMethodNotAllowed(w, r)
		return
	}

	// Read request body.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to read request body: %v", err), http.StatusBadRequest)
		return
	}

	//Parse JSON to DTO
	var task entities.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to unmarshal JSON request body: %v", err), http.StatusBadRequest)
		return
	}

	//Add task to service
	err = h.Service.TaskAdd(r.Context(), task)
	if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (h *TasksHandler) RemoveHandler(w http.ResponseWriter, r *http.Request) {
	// Check HTTP Method
	if r.Method != http.MethodDelete {
		ErrMethodNotAllowed(w, r)
		return
	}

	// Fetching id value
	parts := strings.Split(r.URL.Path, "/")
	urlId := parts[len(parts)-1]

	//Convert it to uint64
	id, err := strconv.ParseUint(urlId, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse id: %v", err), http.StatusBadRequest)
		return
	}

	//Delete task from service
	err = h.Service.TaskRemove(r.Context(), id)

	// If task not exist
	if errors.Is(err, entities.ErrNoTask) {
		http.Error(w, entities.ErrNoTask.Error(), http.StatusNotFound)
		return
	}
	// Other errors
	if err != nil {
		ErrInternalServerError(w, r, http.StatusText(http.StatusNotFound))
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (h *TasksHandler) UpdateHandler(c *fiber.Ctx) error {
	//Read body and parse JSON to DTO
	var task entities.Task
	err := json.Unmarshal(c.Body(), &task)
	if err != nil {
		return fiber.ErrBadRequest
	}

	//Update task in service
	err = h.Service.TaskUpdate(c.Context(), task)
	if errors.Is(err, entities.ErrNoTask) {
		return fiber.ErrNotFound
	}
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return nil
}

func ErrMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	http.Error(w, fmt.Sprintf("Method %s not allowed", r.Method), http.StatusMethodNotAllowed)
}

func ErrInternalServerError(w http.ResponseWriter, r *http.Request, err string) {
	_ = r
	if err == "" {
		err = "Internal server error"
	}
	http.Error(w, err, http.StatusInternalServerError)
}
