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

func (h *TasksHandler) RemoveHandler(w http.ResponseWriter, r *http.Request) {
	// Check HTTP Method
	if r.Method != http.MethodDelete {
		ErrMethodNotAllowed(w, r)
		return
	}

	// Parse URL fields
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse form: %v", err), http.StatusBadRequest)
		return
	}

	// Fetching id value
	urlId := r.FormValue("id")

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
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Other errors
	if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)

}

func ErrMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	http.Error(w, fmt.Sprintf("Method %s not allowed", r.Method), http.StatusMethodNotAllowed)
}

func ErrInternalServerError(w http.ResponseWriter, r *http.Request, err string) {
	if err == "" {
		err = "Internal server error"
	}
	http.Error(w, err, http.StatusInternalServerError)
}
