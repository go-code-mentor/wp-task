package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockedServices struct {
	mock.Mock
}

func (m *MockedServices) GetTask(id uint32) (Task, error) {
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

		s := new(MockedServices)
		s.On("GetTask", task.ID).Return(task, nil)

		result, err := GetTask(s, task.ID)
		s.AssertExpectations(t)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if result != task {
			t.Errorf("expected user: %v, got: %v", task, result)
		}
	})

	t.Run("task getting with error", func(t *testing.T) {
		s := new(MockedServices)
		s.On("GetTask", mock.Anything).Return(Task{}, fmt.Errorf("error"))

		_, err := GetTask(s, 0)
		s.AssertExpectations(t)

		if err == nil {
			t.Errorf("error is nil")
		}
	})
}
