package handlers_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"

	"github.com/go-code-mentor/wp-task/internal/entities"
	"github.com/go-code-mentor/wp-task/internal/handlers"
)

type MockedServices struct {
	mock.Mock
}

func (m *MockedServices) Tasks() ([]entities.Task, error) {
	args := m.Called()
	return args.Get(0).([]entities.Task), args.Error(1)
}

func (m *MockedServices) Task(taskId uint64) (entities.Task, error) {
	args := m.Called(taskId)
	return args.Get(0).(entities.Task), args.Error(1)
}

func TestTaskListHandler(t *testing.T) {

	t.Run("success request", func(t *testing.T) {

		task := entities.Task{
			ID:          5,
			Name:        "test task",
			Description: "test desc",
		}

		app := fiber.New()
		ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

		s := new(MockedServices)

		h := &handlers.TasksHandler{
			Service: s,
		}

		s.On("Tasks").Return([]entities.Task{task}, nil)

		err := h.ListHandler(ctx)
		assert.NoError(t, err)

		result := ctx.Response()

		assert.Equal(t, fasthttp.StatusOK, result.StatusCode())

		body := result.Body()

		encoded := []entities.Task{}
		err = json.Unmarshal(body, &encoded)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(encoded))
		assert.Equal(t, task, encoded[0])

		s.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {

		app := fiber.New()
		ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

		s := new(MockedServices)

		h := &handlers.TasksHandler{
			Service: s,
		}

		s.On("Tasks").Return([]entities.Task{}, fmt.Errorf("error"))

		err := h.ListHandler(ctx)
		assert.Error(t, err)
		assert.Equal(t, err, fiber.ErrInternalServerError)

		s.AssertExpectations(t)
	})
}

func TestTaskItemHandler(t *testing.T) {

	t.Run("success request", func(t *testing.T) {

		// task := entities.Task{
		// 	ID:          "5",
		// 	Name:        "test task",
		// 	Description: "test desc",
		// }

		// app := fiber.New()

		// r := fasthttp.Request{}
		// r.SetRequestURIBytes([]byte("/api/v1/tasks/5"))
		// ctx := app.AcquireCtx(&fasthttp.RequestCtx{
		// 	Request: r,
		// })

		// s := new(MockedServices)

		// h := &handlers.TasksHandler{
		// 	Service: s,
		// }

		// s.On("Task", task.ID).Return(task, nil)

		// err := h.ItemHandler(ctx)
		// assert.NoError(t, err)

		// result := ctx.Response()

		// assert.Equal(t, fasthttp.StatusOK, result.StatusCode())

		// body := result.Body()

		// encoded := entities.Task{}
		// err = json.Unmarshal(body, &encoded)

		// assert.NoError(t, err)
		// assert.Equal(t, task, encoded)

		// s.AssertExpectations(t)
	})

	// t.Run("method not allowed", func(t *testing.T) {
	// 	w := httptest.NewRecorder()

	// 	s := new(MockedServices)
	// 	ctx := context.Background()

	// 	h := &handlers.TasksHandler{
	// 		Service: s,
	// 	}

	// 	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
	// 		r := httptest.NewRequest(method, "/", bytes.NewReader([]byte("")))
	// 		h.ItemHandler(w, r)

	// 		result := w.Result()
	// 		assert.Equal(t, http.StatusMethodNotAllowed, result.StatusCode)

	// 		s.AssertNotCalled(t, "Task", ctx, 0)
	// 		s.AssertExpectations(t)

	// 	}
	// })

	// t.Run("bad request", func(t *testing.T) {

	// 	w := httptest.NewRecorder()
	// 	r := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte("")))

	// 	s := new(MockedServices)
	// 	ctx := context.Background()

	// 	h := &handlers.TasksHandler{
	// 		Service: s,
	// 	}

	// 	s.On("Task", ctx, uint64(0)).Return(entities.Task{}, fmt.Errorf("error"))

	// 	h.ItemHandler(w, r)

	// 	result := w.Result()

	// 	assert.Equal(t, http.StatusBadRequest, result.StatusCode)

	// 	s.AssertExpectations(t)
	// })

	// t.Run("internal server error", func(t *testing.T) {

	// 	w := httptest.NewRecorder()
	// 	r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/users/%d", 0), bytes.NewReader([]byte("")))

	// 	s := new(MockedServices)
	// 	ctx := context.Background()

	// 	h := &handlers.TasksHandler{
	// 		Service: s,
	// 	}

	// 	s.On("Task", ctx, uint64(0)).Return(entities.Task{}, fmt.Errorf("error"))

	// 	h.ItemHandler(w, r)

	// 	result := w.Result()

	// 	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)

	// 	s.AssertExpectations(t)
	// })
}
