package handlers_test

import (
	"context"
	"encoding/json"
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
	return args.Error(1)
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
