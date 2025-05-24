package tgclient

import (
	"context"
	"fmt"

	tgapi "github.com/go-code-mentor/wp-tg-bot/api"
)

type Service struct {
	api tgapi.TgBotClient
}

func New(api tgapi.TgBotClient) *Service {
	return &Service{
		api: api,
	}
}

func (s *Service) SendTask(ctx context.Context, id uint64, name string, description string, login string) error {

	if _, err := s.api.TaskAdd(ctx, &tgapi.TaskAddRequest{
		Id:          id,
		Name:        name,
		Description: description,
		Owner:       login,
	}); err != nil {
		return fmt.Errorf("sending task to tg error: %w", err)
	}

	return nil
}
