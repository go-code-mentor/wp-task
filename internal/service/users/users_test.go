package users_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/go-code-mentor/wp-task/internal/service/users"
)

type MockedStorage struct {
	mock.Mock
}

func (m *MockedStorage) GetUserLogin(ctx context.Context, token string) (string, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(string), args.Error(1)
}

func TestGetUserLogin(t *testing.T) {
	t.Run("success login getting", func(t *testing.T) {
		token := "123"
		login := "user"
		ctx := context.Background()
		storageMock := new(MockedStorage)
		storageMock.On("GetUserLogin", ctx, token).Return(login, nil)
		s := users.New(storageMock)

		result, err := s.GetUserLogin(ctx, token)
		assert.NoError(t, err)
		assert.Equal(t, login, result)
	})

	t.Run("login getting with error", func(t *testing.T) {
		token := "123"
		ctx := context.Background()
		storageMock := new(MockedStorage)
		storageMock.On("GetUserLogin", ctx, token).Return("", fmt.Errorf("error"))
		s := users.New(storageMock)

		_, err := s.GetUserLogin(ctx, token)
		assert.Error(t, err)
	})

}
