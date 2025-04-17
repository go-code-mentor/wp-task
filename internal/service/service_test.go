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

type MockedStorage struct {
	mock.Mock
}

func (m *MockedStorage) Task(ctx context.Context, id uint64) (entities.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entities.Task), args.Error(1)
}

func (m *MockedStorage) Tasks(ctx context.Context) ([]entities.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entities.Task), args.Error(1)
}

func (m *MockedStorage) TaskRemove(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
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
		storageMock := new(MockedStorage)
		storageMock.On("Task", ctx, task.ID).Return(task, nil)
		s := service.New(ctx, storageMock)

		result, err := s.Task(task.ID)
		assert.NoError(t, err)
		assert.Equal(t, task, result)
	})

	t.Run("task getting with error", func(t *testing.T) {
		taskId := uint64(1)
		ctx := context.Background()
		storageMock := new(MockedStorage)
		storageMock.On("Task", ctx, taskId).Return(entities.Task{}, fmt.Errorf("error"))
		s := service.New(ctx, storageMock)

		_, err := s.Task(taskId)
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
		storageMock := new(MockedStorage)
		storageMock.On("Tasks", ctx).Return([]entities.Task{task}, nil)
		s := service.New(ctx, storageMock)

		result, err := s.Tasks()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, task, result[0])
	})

	t.Run("tasks getting with error", func(t *testing.T) {
		ctx := context.Background()
		storageMock := new(MockedStorage)
		storageMock.On("Tasks", ctx).Return([]entities.Task{}, fmt.Errorf("error"))
		s := service.New(ctx, storageMock)

		_, err := s.Tasks()
		assert.Error(t, err)
	})

}

func TestTaskRemoving(t *testing.T) {
	t.Run("success task removing", func(t *testing.T) {
		taskId := uint64(1)
		ctx := context.Background()
		storageMock := new(MockedStorage)
		storageMock.On("TaskRemove", ctx, taskId).Return(nil)
		s := service.New(ctx, storageMock)

		err := s.TaskRemove(taskId)
		assert.NoError(t, err)
	})

	t.Run("task removing with error", func(t *testing.T) {
		taskId := uint64(1)
		ctx := context.Background()
		storageMock := new(MockedStorage)
		storageMock.On("TaskRemove", ctx, taskId).Return(fmt.Errorf("error"))
		s := service.New(ctx, storageMock)

		err := s.TaskRemove(taskId)
		assert.Error(t, err)
	})
}
