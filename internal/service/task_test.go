package service_test

import (
	"context"
	"fmt"
	"github.com/go-code-mentor/wp-task/internal/entities"
	"github.com/go-code-mentor/wp-task/internal/service"
	"github.com/stretchr/testify/assert"
	"testing"
)

func (m *MockedServices) TaskUpdate(ctx context.Context, task entities.Task) (entities.Task, error) {
	args := m.Called(task)
	return args.Get(0).(entities.Task), args.Error(1)

}

func TestTaskUpdating(t *testing.T) {
	t.Run("success task updating", func(t *testing.T) {
		task := entities.Task{
			ID:          1,
			Name:        "Test task",
			Description: "test task description",
		}
		ctx := context.Background()
		storageMock := new(MockedServices)
		storageMock.On("TaskUpdate", task).Return(task, nil)
		s := service.New(storageMock)

		result, err := s.TaskUpdate(ctx, task)
		assert.NoError(t, err)
		assert.Equal(t, task, result)
	})

	t.Run("task updating with error", func(t *testing.T) {
		task := entities.Task{
			ID:          1,
			Name:        "Test task",
			Description: "test task description",
		}
		ctx := context.Background()
		storageMock := new(MockedServices)
		storageMock.On("TaskUpdate", task).Return(entities.Task{}, fmt.Errorf("error"))
		s := service.New(storageMock)

		_, err := s.TaskUpdate(ctx, task)
		assert.Error(t, err)

	})

}
