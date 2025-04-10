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

func (m *MockedServices) Tasks(ctx context.Context) ([]entities.Task, error) {
	args := m.Called()
	return args.Get(0).([]entities.Task), args.Error(1)
}

func (m *MockedServices) TaskRemove(ctx context.Context, id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestTaskGetting(t *testing.T) {
	t.Run("success task getting", func(t *testing.T) {
		task := entities.Task{
			ID:          1,
			Name:        "Test task",
			Description: "test task description",
		}
		ctx := context.Background()
		storageMock := new(MockedServices)
		storageMock.On("Task", task.ID).Return(task, nil)
		s := service.New(storageMock)

		result, err := s.Task(ctx, task.ID)
		assert.NoError(t, err)
		assert.Equal(t, task, result)
	})

	t.Run("task getting with error", func(t *testing.T) {
		taskId := uint64(1)
		ctx := context.Background()
		storageMock := new(MockedServices)
		storageMock.On("Task", taskId).Return(entities.Task{}, fmt.Errorf("error"))
		s := service.New(storageMock)

		_, err := s.Task(ctx, taskId)
		assert.Error(t, err)
	})
}

func TestTasksGetting(t *testing.T) {
	t.Run("success tasks getting", func(t *testing.T) {
		task := entities.Task{
			ID:          1,
			Name:        "Test task",
			Description: "test task description",
		}
		ctx := context.Background()
		storageMock := new(MockedServices)
		storageMock.On("Tasks").Return([]entities.Task{task}, nil)
		s := service.New(storageMock)

		result, err := s.Tasks(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, task, result[0])
	})

	t.Run("tasks getting with error", func(t *testing.T) {
		ctx := context.Background()
		storageMock := new(MockedServices)
		storageMock.On("Tasks").Return([]entities.Task{}, fmt.Errorf("error"))
		s := service.New(storageMock)

		_, err := s.Tasks(ctx)
		assert.Error(t, err)
	})

}

func TestTaskRemoving(t *testing.T) {
	t.Run("success task removing", func(t *testing.T) {
		taskId := uint64(1)
		ctx := context.Background()
		storageMock := new(MockedServices)
		storageMock.On("TaskRemove", taskId).Return(nil)
		s := service.New(storageMock)

		err := s.TaskRemove(ctx, taskId)
		assert.NoError(t, err)
	})

	t.Run("task removing with error", func(t *testing.T) {
		taskId := uint64(1)
		ctx := context.Background()
		storageMock := new(MockedServices)
		storageMock.On("TaskRemove", taskId).Return(fmt.Errorf("error"))
		s := service.New(storageMock)

		err := s.TaskRemove(ctx, taskId)
		assert.Error(t, err)
	})
}
