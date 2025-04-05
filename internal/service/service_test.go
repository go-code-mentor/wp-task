package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/go-code-mentor/wp-task/internal/entities"
	"github.com/go-code-mentor/wp-task/internal/service"
)

type MockedServices struct {
	mock.Mock
}

func (m *MockedServices) Task(ctx context.Context, id uint64) (entities.Task, error) {
	args := m.Called(id)
	return args.Get(0).(entities.Task), args.Error(1)
}

func TestTaskGetting(t *testing.T) {

	t.Run("success task getting", func(t *testing.T) {

		task := entities.Task{
			ID:          1,
			Name:        "Test task",
			Description: "test task description",
		}

		ctx := context.Background()
		s := new(MockedServices)
		s.On("Task", task.ID).Return(task, nil)

		result, err := service.Task(ctx, s, task.ID)

		s.AssertExpectations(t)

		assert.NoError(t, err)

		assert.Equal(t, task, result)

	})

	t.Run("task getting with error", func(t *testing.T) {
		taskId := uint64(1)
		ctx := context.Background()
		s := new(MockedServices)
		s.On("Task", taskId).Return(entities.Task{}, fmt.Errorf("error"))

		_, err := service.Task(ctx, s, taskId)
		s.AssertExpectations(t)

		assert.Error(t, err)
	})
}
