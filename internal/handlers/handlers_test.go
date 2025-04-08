package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestTaskListHandler(t *testing.T) {

	t.Run("success request", func(t *testing.T) {

		task := entities.Task{
			ID:          5,
			Name:        "test task",
			Description: "test desc",
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte("")))

		s := new(MockedServices)
		ctx := context.Background()

		h := &handlers.TasksHandler{
			Service: s,
		}

		s.On("Tasks", ctx).Return([]entities.Task{task}, nil)

		h.ListHandler(w, r)

		result := w.Result()

		assert.Equal(t, http.StatusOK, result.StatusCode)

		body, err := io.ReadAll(result.Body)
		if err != nil {
			t.Fatal("Error reading response body:", err)
		}

		if result.Body != nil {
			result.Body.Close()
		}

		encoded := []entities.Task{}
		err = json.Unmarshal(body, &encoded)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(encoded))
		assert.Equal(t, task, encoded[0])

		s.AssertExpectations(t)
	})

	t.Run("method not allowed", func(t *testing.T) {
		w := httptest.NewRecorder()

		s := new(MockedServices)
		ctx := context.Background()

		h := &handlers.TasksHandler{
			Service: s,
		}

		for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
			r := httptest.NewRequest(method, "/", bytes.NewReader([]byte("")))
			h.ListHandler(w, r)

			result := w.Result()
			assert.Equal(t, http.StatusMethodNotAllowed, result.StatusCode)

			s.AssertNotCalled(t, "Tasks", ctx)
			s.AssertExpectations(t)

		}
	})

	t.Run("internal server error", func(t *testing.T) {

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte("")))

		s := new(MockedServices)
		ctx := context.Background()

		h := &handlers.TasksHandler{
			Service: s,
		}

		s.On("Tasks", ctx).Return([]entities.Task{}, fmt.Errorf("error"))

		h.ListHandler(w, r)

		result := w.Result()

		assert.Equal(t, http.StatusInternalServerError, result.StatusCode)

		s.AssertExpectations(t)
	})
}
