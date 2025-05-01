package simpletoken

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"github.com/go-code-mentor/wp-task/internal/entities"
)

const (
	AuthHeader = "Authorization"
)

type AuthService interface {
	GetUserLogin(ctx context.Context, token string) (string, error)
}

type AuthMiddleware struct {
	Service AuthService
}

func (m *AuthMiddleware) Auth(c *fiber.Ctx) error {

	token := c.Get(AuthHeader, "")

	if len(token) == 0 {
		return fiber.ErrUnauthorized
	}

	userLogin, err := m.Service.GetUserLogin(c.Context(), token)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	c.Locals(entities.UserLoginKey, userLogin)

	return c.Next()
}
