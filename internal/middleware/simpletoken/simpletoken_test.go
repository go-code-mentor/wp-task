package simpletoken_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/go-code-mentor/wp-task/internal/entities"
	"github.com/go-code-mentor/wp-task/internal/middleware/simpletoken"
)

type MockedUserService struct {
	mock.Mock
}

func (m *MockedUserService) GetUserLogin(ctx context.Context, token string) (string, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(string), args.Error(1)
}

func TestAuthMiddleware(t *testing.T) {
	t.Run("success auth", func(t *testing.T) {
		token := "123"
		login := "user"

		serviceMock := new(MockedUserService)
		serviceMock.On("GetUserLogin", mock.Anything, token).Return(login, nil)
		authMiddleware := simpletoken.AuthMiddleware{Service: serviceMock}

		app := fiber.New()
		app.Use(authMiddleware.Auth)
		app.Get("/dummy", func(c *fiber.Ctx) error {
			l, ok := c.Locals(entities.UserLoginKey).(string)
			if !ok {
				t.Errorf("getting login from context failed: want: %v got: %v", login, l)
			}
			assert.Equal(t, l, login)
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/dummy", nil)
		req.Header.Set(simpletoken.AuthHeader, token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("unauthorized error if auth header not set", func(t *testing.T) {

		serviceMock := new(MockedUserService)
		authMiddleware := simpletoken.AuthMiddleware{Service: serviceMock}

		app := fiber.New()
		app.Use(authMiddleware.Auth)
		app.Get("/dummy", nil)

		req := httptest.NewRequest(http.MethodGet, "/dummy", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		serviceMock.AssertNotCalled(t, "GetUserLogin", mock.Anything)
		serviceMock.AssertExpectations(t)
	})

	t.Run("unauthorized error if auth header is empty", func(t *testing.T) {

		serviceMock := new(MockedUserService)
		authMiddleware := simpletoken.AuthMiddleware{Service: serviceMock}

		app := fiber.New()
		app.Use(authMiddleware.Auth)
		app.Get("/dummy", nil)

		req := httptest.NewRequest(http.MethodGet, "/dummy", nil)
		req.Header.Set(simpletoken.AuthHeader, "")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		serviceMock.AssertNotCalled(t, "GetUserLogin", mock.Anything)
		serviceMock.AssertExpectations(t)
	})

	t.Run("unauthorized error if auth token not exists", func(t *testing.T) {

		token := "123"

		serviceMock := new(MockedUserService)
		serviceMock.On("GetUserLogin", mock.Anything, token).Return("", fmt.Errorf("error"))
		authMiddleware := simpletoken.AuthMiddleware{Service: serviceMock}

		app := fiber.New()
		app.Use(authMiddleware.Auth)
		app.Get("/dummy", nil)

		req := httptest.NewRequest(http.MethodGet, "/dummy", nil)
		req.Header.Set(simpletoken.AuthHeader, token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		serviceMock.AssertExpectations(t)
	})

}
