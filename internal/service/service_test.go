package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedServices struct {
	mock.Mock
}

func (m *MockedServices) Task(id uint32) (Task, error) {
	args := m.Called(id)
	return args.Get(0).(Task), args.Error(1)
}

func TestTaskGetting(t *testing.T) {

	t.Run("success task getting", func(t *testing.T) {

		task := Task{
			ID:          1,
			Name:        "Test task",
			Description: "test task description",
		}

		ctx := context.Background()
		s := new(MockedServices)
		s.On("Task", task.ID).Return(task, nil)

		result, err := GetTask(ctx, s, task.ID)

		s.AssertExpectations(t)

		assert.NoError(t, err)

		assert.Equal(t, task, result)

	})

	t.Run("task getting with error", func(t *testing.T) {
		taskId := uint32(1)
		ctx := context.Background()
		s := new(MockedServices)
		s.On("Task", taskId).Return(Task{}, fmt.Errorf("error"))

		_, err := GetTask(ctx, s, taskId)
		s.AssertExpectations(t)

		assert.Error(t, err)
	})
}
