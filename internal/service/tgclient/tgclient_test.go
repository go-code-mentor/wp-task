package tgclient_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	"github.com/go-code-mentor/wp-task/internal/service/tgclient"
	tgapi "github.com/go-code-mentor/wp-tg-bot/api"
)

type MockedApi struct {
	mock.Mock
}

func (m *MockedApi) TaskAdd(ctx context.Context, in *tgapi.TaskAddRequest, opts ...grpc.CallOption) (*tgapi.TaskAddResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*tgapi.TaskAddResponse), args.Error(1)
}

func (m *MockedApi) Ping(ctx context.Context, in *tgapi.PingRequest, opts ...grpc.CallOption) (*tgapi.PingResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*tgapi.PingResponse), args.Error(1)
}

func TestSendTask(t *testing.T) {
	t.Run("success task sending", func(t *testing.T) {
		task := &tgapi.TaskAddRequest{
			Id:          1,
			Name:        "Test task",
			Description: "test task description",
			Owner:       "user",
		}
		ctx := context.Background()
		apiMock := new(MockedApi)
		apiMock.On("TaskAdd", ctx, task).Return(&tgapi.TaskAddResponse{}, nil)
		service := tgclient.New(apiMock)

		err := service.SendTask(ctx, task.Id, task.Name, task.Description, task.Owner)
		assert.NoError(t, err)
	})

	t.Run("task sending with error", func(t *testing.T) {
		task := &tgapi.TaskAddRequest{
			Id:          1,
			Name:        "Test task",
			Description: "test task description",
			Owner:       "user",
		}
		ctx := context.Background()
		apiMock := new(MockedApi)
		apiMock.On("TaskAdd", ctx, task).Return(&tgapi.TaskAddResponse{}, fmt.Errorf("error"))
		service := tgclient.New(apiMock)

		err := service.SendTask(ctx, task.Id, task.Name, task.Description, task.Owner)
		assert.Error(t, err)
	})
}
