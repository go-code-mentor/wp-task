package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-code-mentor/wp-task/internal/entities"
	"net/http"
)

type Service interface {
	Tasks(ctx context.Context) ([]entities.Task, error)
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

func ErrMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	http.Error(w, fmt.Sprintf("Method %s not allowed", r.Method), http.StatusMethodNotAllowed)
}

func ErrInternalServerError(w http.ResponseWriter, r *http.Request, err string) {
	if err == "" {
		err = "Internal server error"
	}
	http.Error(w, err, http.StatusInternalServerError)
}
