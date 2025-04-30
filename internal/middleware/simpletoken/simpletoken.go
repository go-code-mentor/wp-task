package simpletoken

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

const (
	AuthHeader   = "Authorization"
	UserLoginKey = "user"
)

type AuthService interface {
	Auth(ctx context.Context, token string) (string, error)
}

type AuthMiddleware struct {
	Service AuthService
}

func (m *AuthMiddleware) Auth(c *fiber.Ctx) error {

	token := c.Get(fiber.HeaderAuthorization, "")

	if len(token) == 0 {
		return fiber.ErrUnauthorized
	}

	userLogin, err := m.Service.Auth(c.Context(), token)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	c.Locals(UserLoginKey, userLogin)

	return c.Next()
}
