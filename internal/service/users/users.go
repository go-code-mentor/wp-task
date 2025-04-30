package users

import (
	"context"
	"fmt"
)

type UserStorage interface {
	Auth(ctx context.Context, token string) (string, error)
}

func New(storage UserStorage) *UserService {
	return &UserService{
		Storage: storage,
	}
}

type UserService struct {
	Storage UserStorage
}

func (s *UserService) Auth(ctx context.Context, token string) (string, error) {
	login, err := s.Storage.Auth(ctx, token)
	if err != nil {
		return login, fmt.Errorf("could not auth user: %w", err)
	}
	return login, nil
}
