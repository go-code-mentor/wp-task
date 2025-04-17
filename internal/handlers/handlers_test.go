package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/go-code-mentor/wp-task/internal/entities"
	"github.com/go-code-mentor/wp-task/internal/handlers"
)

type MockedServices struct {
	mock.Mock
}

func (m *MockedServices) Tasks(ctx context.Context) ([]entities.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entities.Task), args.Error(1)
}

func (m *MockedServices) Task(ctx context.Context, taskId uint64) (entities.Task, error) {
	args := m.Called(ctx, taskId)
	return args.Get(0).(entities.Task), args.Error(1)
}

func (m *MockedServices) TaskAdd(ctx context.Context, task entities.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockedServices) TaskRemove(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestTaskListHandler(t *testing.T) {

	t.Run("success request", func(t *testing.T) {

		task := entities.Task{
			ID:          5,
			Name:        "test task",
			Description: "test desc",
		}

		s := new(MockedServices)
		h := &handlers.TasksHandler{
			Service: s,
		}
		s.On("Tasks", mock.Anything).Return([]entities.Task{task}, nil)

		app := fiber.New()
		app.Get("/tasks", h.ListHandler)

		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body := make([]byte, resp.ContentLength)

		n, err := resp.Body.Read(body)
		if n != int(resp.ContentLength) {
			t.Fatal("Error reading response body:", err)
		}
		if err != nil && err != io.EOF {
			t.Fatal("Error reading response body:", err)
		}

		err = resp.Body.Close()
		if err != nil {
			t.Fatal("Error closing body:", err)
		}

		encoded := []entities.Task{}
		err = json.Unmarshal(body, &encoded)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(encoded))
		assert.Equal(t, task, encoded[0])

		// s.AssertExpectations(t)
	})

	t.Run("method not allowed", func(t *testing.T) {

		s := new(MockedServices)
		h := &handlers.TasksHandler{
			Service: s,
		}

		app := fiber.New()
		app.Get("/tasks", h.ListHandler)

		for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {

			req := httptest.NewRequest(method, "/tasks", nil)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

			s.AssertNotCalled(t, "Tasks", mock.Anything)
			s.AssertExpectations(t)
		}
	})

	t.Run("internal server error", func(t *testing.T) {

		s := new(MockedServices)
		h := &handlers.TasksHandler{
			Service: s,
		}
		s.On("Tasks", mock.Anything).Return([]entities.Task{}, fmt.Errorf("error"))

		app := fiber.New()
		app.Get("/tasks", h.ListHandler)

		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		s.AssertExpectations(t)
	})
}

func TestTaskItemHandler(t *testing.T) {

	t.Run("success request", func(t *testing.T) {

		task := entities.Task{
			ID:          5,
			Name:        "test task",
			Description: "test desc",
		}

		s := new(MockedServices)
		h := &handlers.TasksHandler{
			Service: s,
		}
		s.On("Task", mock.Anything, task.ID).Return(task, nil)

		app := fiber.New()
		app.Get("/tasks/:id", h.ItemHandler)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/tasks/%d", task.ID), nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body := make([]byte, resp.ContentLength)
		n, err := resp.Body.Read(body)
		if n != int(resp.ContentLength) {
			t.Fatal("Error reading response body:", err)
		}
		if err != nil && err != io.EOF {
			t.Fatal("Error reading response body:", err)
		}

		err = resp.Body.Close()
		if err != nil {
			t.Fatal("Error closing body:", err)
		}

		encoded := entities.Task{}
		err = json.Unmarshal(body, &encoded)

		assert.NoError(t, err)
		assert.Equal(t, task, encoded)

		s.AssertExpectations(t)
	})

	t.Run("method not allowed", func(t *testing.T) {

		taskId := uint64(1)

		s := new(MockedServices)
		h := &handlers.TasksHandler{
			Service: s,
		}

		app := fiber.New()
		app.Get("/tasks/:id", h.ItemHandler)

		for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
			req := httptest.NewRequest(method, fmt.Sprintf("/tasks/%d", taskId), nil)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

			s.AssertNotCalled(t, "Task", mock.Anything, taskId)
			s.AssertExpectations(t)
		}
	})

	t.Run("error not found - wrong type of task id", func(t *testing.T) {

		taskId := "abc"

		s := new(MockedServices)
		h := &handlers.TasksHandler{
			Service: s,
		}

		app := fiber.New()
		app.Get("/tasks/:id", h.ItemHandler)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/tasks/%s", taskId), nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		s.AssertNotCalled(t, "Task", mock.Anything, taskId)
		s.AssertExpectations(t)
	})

	t.Run("error not found - task not exists in storage", func(t *testing.T) {

		taskId := uint64(1)

		s := new(MockedServices)
		h := &handlers.TasksHandler{
			Service: s,
		}
		s.On("Task", mock.Anything, taskId).Return(entities.Task{}, fmt.Errorf("error"))

		app := fiber.New()
		app.Get("/tasks/:id", h.ItemHandler)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/tasks/%d", taskId), nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		s.AssertExpectations(t)
	})
}

func TestTaskAddHandler(t *testing.T) {
	t.Run("success request", func(t *testing.T) {
		task := entities.Task{
			ID:          1,
			Name:        "Test task",
			Description: "test task description",
		}

		requestBody, err := json.Marshal(task)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBody))

		s := new(MockedServices)
		ctx := context.Background()

		h := &handlers.TasksHandler{
			Service: s,
		}

		s.On("TaskAdd", ctx, task).Return(nil)

		h.AddHandler(w, r)

		result := w.Result()

		assert.Equal(t, http.StatusOK, result.StatusCode)
		s.AssertExpectations(t)
	})

	t.Run("method not allowed", func(t *testing.T) {
		w := httptest.NewRecorder()
		s := new(MockedServices)
		h := &handlers.TasksHandler{Service: s}

		// Проверяем все неподходящие методы
		for _, method := range []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodDelete} {
			r := httptest.NewRequest(method, "/", bytes.NewReader([]byte("")))
			h.AddHandler(w, r)

			result := w.Result()

			assert.Equal(t, http.StatusMethodNotAllowed, result.StatusCode)
			s.AssertNotCalled(t, "TaskAdd", mock.Anything, mock.Anything)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("{invalid json}")))

		s := new(MockedServices)
		h := &handlers.TasksHandler{Service: s}

		h.AddHandler(w, r)

		result := w.Result()

		assert.Equal(t, http.StatusBadRequest, result.StatusCode)
		s.AssertNotCalled(t, "TaskAdd", mock.Anything, mock.Anything)
	})

	t.Run("internal server error", func(t *testing.T) {
		task := entities.Task{
			ID:          1,
			Name:        "Test task",
			Description: "test task description",
		}

		requestBody, err := json.Marshal(task)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBody))

		s := new(MockedServices)
		ctx := context.Background()

		h := &handlers.TasksHandler{
			Service: s,
		}

		s.On("TaskAdd", ctx, task).Return(errors.New("service error"))

		h.AddHandler(w, r)

		result := w.Result()

		assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
		s.AssertExpectations(t)
	})
}

func TestTaskRemoveHandler(t *testing.T) {
	t.Run("success request", func(t *testing.T) {
		var id uint64
		id = 1

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/task/%d", id), nil)

		s := new(MockedServices)
		ctx := context.Background()

		h := &handlers.TasksHandler{
			Service: s,
		}

		s.On("TaskRemove", ctx, id).Return(nil)

		h.RemoveHandler(w, r)

		result := w.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusOK, result.StatusCode)
		s.AssertExpectations(t)
	})

	t.Run("method not allowed", func(t *testing.T) {
		var id uint64
		id = 1
		w := httptest.NewRecorder()

		s := new(MockedServices)
		ctx := context.Background()

		h := &handlers.TasksHandler{
			Service: s,
		}

		for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodGet} {
			r := httptest.NewRequest(method, fmt.Sprintf("/api/v1/task/%d", id), nil)

			h.RemoveHandler(w, r)

			result := w.Result()
			assert.Equal(t, http.StatusMethodNotAllowed, result.StatusCode)

			s.AssertNotCalled(t, "TaskRemove", ctx, id)
			s.AssertExpectations(t)

		}
	})

	t.Run("invalid id format (not uint64)", func(t *testing.T) {
		s := new(MockedServices)
		ctx := context.Background()

		h := &handlers.TasksHandler{
			Service: s,
		}

		for _, id := range []string{"string", "-1", "18446744073709551616", ""} {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/task/%v", id), nil)

			h.RemoveHandler(w, r)

			result := w.Result()
			assert.Equal(t, http.StatusBadRequest, result.StatusCode)

			s.AssertNotCalled(t, "TaskRemove", ctx, id)
			s.AssertExpectations(t)

		}
	})

	t.Run("task not found", func(t *testing.T) {
		var id uint64
		id = 1

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/task/%d", id), nil)

		s := new(MockedServices)
		ctx := context.Background()

		h := &handlers.TasksHandler{
			Service: s,
		}

		s.On("TaskRemove", ctx, id).Return(entities.ErrNoTask)

		h.RemoveHandler(w, r)

		result := w.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusNotFound, result.StatusCode)
		s.AssertExpectations(t)
	})

	t.Run("task removing with error", func(t *testing.T) {
		var id uint64
		id = 1

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/task/%d", id), nil)

		s := new(MockedServices)
		ctx := context.Background()

		h := &handlers.TasksHandler{
			Service: s,
		}

		s.On("TaskRemove", ctx, id).Return(fmt.Errorf("error"))

		h.RemoveHandler(w, r)

		result := w.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
		s.AssertExpectations(t)
	})
}
