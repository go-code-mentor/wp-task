package service_test

import (
	"context"
	"fmt"
	"github.com/go-code-mentor/wp-task/internal/entities"
	"github.com/go-code-mentor/wp-task/internal/service"
	"github.com/stretchr/testify/assert"
	"testing"
)

func (m *MockedServices) TaskUpdate(ctx context.Context, task entities.Task) error {
	args := m.Called(task)
	return args.Error(0)

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
		storageMock.On("TaskUpdate", task).Return(nil)
		s := service.New(storageMock)

		err := s.TaskUpdate(ctx, task)
		assert.NoError(t, err)
	})

	t.Run("task updating with error", func(t *testing.T) {
		task := entities.Task{
			ID:          1,
			Name:        "Test task",
			Description: "test task description",
		}
		ctx := context.Background()
		storageMock := new(MockedServices)
		storageMock.On("TaskUpdate", task).Return(fmt.Errorf("error"))
		s := service.New(storageMock)

		err := s.TaskUpdate(ctx, task)
		assert.Error(t, err)

	})

}
