package tgclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	tgapi "github.com/go-code-mentor/wp-tg-bot/api"
)

type Service struct {
	api tgapi.TgBotClient
}

func New(conn *grpc.ClientConn) *Service {
	return &Service{
		api: tgapi.NewTgBotClient(conn),
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
