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

func (m *MockedStorage) Task(ctx context.Context, id uint64, login string) (entities.Task, error) {
	args := m.Called(ctx, id, login)
	return args.Get(0).(entities.Task), args.Error(1)
}

func (m *MockedStorage) Tasks(ctx context.Context, login string) ([]entities.Task, error) {
	args := m.Called(ctx, login)
	return args.Get(0).([]entities.Task), args.Error(1)
}

func (m *MockedStorage) TaskRemove(ctx context.Context, id uint64, login string) error {
	args := m.Called(ctx, id, login)
	return args.Error(0)
}

func (m *MockedStorage) TaskUpdate(ctx context.Context, task entities.Task, login string) error {
	args := m.Called(ctx, task, login)
	return args.Error(0)

}

func (m *MockedStorage) TaskAdd(ctx context.Context, task entities.Task, login string) (uint64, error) {
	args := m.Called(ctx, task, login)
	return args.Get(0).(uint64), args.Error(1)
}

type MockedTgClient struct {
	mock.Mock
}

func (m *MockedTgClient) SendTask(ctx context.Context, id uint64, name string, description string, login string) error {
	args := m.Called(ctx, id, name, description, login)
	return args.Error(0)
}

func TestTaskGetting(t *testing.T) {
	t.Run("success task getting", func(t *testing.T) {
		task := entities.Task{
			ID:          1,
			Name:        "Test task",
			Description: "test task description",
			Owner:       "user",
		}
		ctx := context.Background()
		storageMock := new(MockedStorage)
		storageMock.On("Task", ctx, task.ID, task.Owner).Return(task, nil)
		s := service.New(storageMock, new(MockedTgClient))

		result, err := s.Task(ctx, task.ID, task.Owner)
		assert.NoError(t, err)
		assert.Equal(t, task, result)
	})

	t.Run("task getting with error", func(t *testing.T) {
		taskId := uint64(1)
		login := "user"
		ctx := context.Background()
		storageMock := new(MockedStorage)
		storageMock.On("Task", ctx, taskId, login).Return(entities.Task{}, fmt.Errorf("error"))
		s := service.New(storageMock, new(MockedTgClient))

		_, err := s.Task(ctx, taskId, login)
		assert.Error(t, err)
	})
}

func TestTasksGetting(t *testing.T) {
	t.Run("success tasks getting", func(t *testing.T) {
		task := entities.Task{
			ID:          1,
			Name:        "Test task",
			Description: "test task description",
			Owner:       "user",
		}
		ctx := context.Background()
		storageMock := new(MockedStorage)
		storageMock.On("Tasks", ctx, task.Owner).Return([]entities.Task{task}, nil)
		s := service.New(storageMock, new(MockedTgClient))

		result, err := s.Tasks(ctx, task.Owner)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, task, result[0])
	})

	t.Run("tasks getting with error", func(t *testing.T) {
		login := "user"
		ctx := context.Background()
		storageMock := new(MockedStorage)
		storageMock.On("Tasks", ctx, login).Return([]entities.Task{}, fmt.Errorf("error"))
		s := service.New(storageMock, new(MockedTgClient))

		_, err := s.Tasks(ctx, login)
		assert.Error(t, err)
	})

}

func TestTaskRemoving(t *testing.T) {
	t.Run("success task removing", func(t *testing.T) {
		taskId := uint64(1)
		login := "user"
		ctx := context.Background()
		storageMock := new(MockedStorage)
		storageMock.On("TaskRemove", ctx, taskId, login).Return(nil)
		s := service.New(storageMock, new(MockedTgClient))

		err := s.TaskRemove(ctx, taskId, login)
		assert.NoError(t, err)
	})

	t.Run("task removing with error", func(t *testing.T) {
		taskId := uint64(1)
		login := "user"
		ctx := context.Background()
		storageMock := new(MockedStorage)
		storageMock.On("TaskRemove", ctx, taskId, login).Return(fmt.Errorf("error"))
		s := service.New(storageMock, new(MockedTgClient))

		err := s.TaskRemove(ctx, taskId, login)
		assert.Error(t, err)
	})
}

func TestTaskAdding(t *testing.T) {
	t.Run("success task adding", func(t *testing.T) {
		task := entities.Task{
			Name:        "Test task",
			Description: "test task description",
			Owner:       "user",
		}
		taskID := uint64(1)
		ctx := context.Background()
		storageMock := new(MockedStorage)
		storageMock.On("TaskAdd", ctx, task, task.Owner).Return(taskID, nil)
		tgClientMock := new(MockedTgClient)
		tgClientMock.On("SendTask", ctx, taskID, task.Name, task.Description, task.Owner).Return(nil)
		s := service.New(storageMock, tgClientMock)

		id, err := s.TaskAdd(ctx, task, task.Owner)
		assert.NoError(t, err)
		assert.Equal(t, taskID, id)
	})

	t.Run("task adding with storage service error", func(t *testing.T) {
		taskId := uint64(1)
		task := entities.Task{
			Name:        "Test task",
			Description: "test task description",
			Owner:       "user",
		}
		ctx := context.Background()
		storageMock := new(MockedStorage)
		storageMock.On("TaskAdd", ctx, task, task.Owner).Return(taskId, fmt.Errorf("error"))
		tgClientMock := new(MockedTgClient)
		tgClientMock.On("SendTask", ctx, taskId, task.Name, task.Description, task.Owner).Return(nil)
		s := service.New(storageMock, new(MockedTgClient))

		id, err := s.TaskAdd(ctx, task, task.Owner)
		assert.Error(t, err)
		assert.Equal(t, uint64(0), id)
		tgClientMock.AssertNotCalled(t, "SendTask", mock.Anything)
	})

	t.Run("task adding with tg client service error", func(t *testing.T) {
		taskId := uint64(1)
		task := entities.Task{
			Name:        "Test task",
			Description: "test task description",
			Owner:       "user",
		}
		ctx := context.Background()
		storageMock := new(MockedStorage)
		storageMock.On("TaskAdd", ctx, task, task.Owner).Return(taskId, nil)
		tgClientMock := new(MockedTgClient)
		tgClientMock.On("SendTask", ctx, taskId, task.Name, task.Description, task.Owner).Return(fmt.Errorf("error"))
		s := service.New(storageMock, tgClientMock)

		id, err := s.TaskAdd(ctx, task, task.Owner)
		assert.Error(t, err)
		assert.Equal(t, uint64(0), id)
	})

}

func TestTaskUpdating(t *testing.T) {
	t.Run("success task updating", func(t *testing.T) {
		task := entities.Task{
			ID:          1,
			Name:        "Test task",
			Description: "test task description",
			Owner:       "user",
		}
		ctx := context.Background()
		storageMock := new(MockedStorage)
		storageMock.On("TaskUpdate", ctx, task, task.Owner).Return(nil)
		s := service.New(storageMock, new(MockedTgClient))

		err := s.TaskUpdate(ctx, task, task.Owner)
		assert.NoError(t, err)
	})

	t.Run("task updating with error", func(t *testing.T) {
		task := entities.Task{
			ID:          1,
			Name:        "Test task",
			Description: "test task description",
			Owner:       "user",
		}
		ctx := context.Background()
		storageMock := new(MockedStorage)
		storageMock.On("TaskUpdate", ctx, task, task.Owner).Return(fmt.Errorf("error"))
		s := service.New(storageMock, new(MockedTgClient))

		err := s.TaskUpdate(ctx, task, task.Owner)
		assert.Error(t, err)
	})
}
