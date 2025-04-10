package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-code-mentor/wp-task/internal/entities"
)

type Service interface {
	Tasks(ctx context.Context) ([]entities.Task, error)
	Task(ctx context.Context, id uint64) (entities.Task, error)
}

type TasksHandler struct {
	Service Service
}

func (h *TasksHandler) ListHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		ErrMethodNotAllowed(w, r)
		return
	}

	tasks, err := h.Service.Tasks(r.Context())
	if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	response, err := json.Marshal(tasks)
	if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (h *TasksHandler) GetItemHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		ErrMethodNotAllowed(w, r)
		return
	}

	userId, err := GetUserId(r.URL.Path)
	if err != nil {
		ErrBadRequest(w, r, "")
	}

	task, err := h.Service.Task(r.Context(), userId)
	if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	response, err := json.Marshal(task)
	if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func GetUserId(path string) (uint64, error) {

	if len(path) < 11 {
		return 0, fmt.Errorf("unexpected users endpoint path. It must to starts with \"/api/users/\"")
	}

	id, err := strconv.Atoi(path[len("/api/users/"):])
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
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

func ErrBadRequest(w http.ResponseWriter, r *http.Request, err string) {
	if err == "" {
		err = "Bad request"
	}
	http.Error(w, err, http.StatusBadRequest)
}
