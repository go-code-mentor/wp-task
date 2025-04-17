package service

import (
	"context"
	"fmt"
	"github.com/go-code-mentor/wp-task/internal/entities"
)

func (s *Service) TaskUpdate(ctx context.Context, task entities.Task) error {
	err := s.Storage.TaskUpdate(ctx, task)
	if err != nil {
		return fmt.Errorf("unable to update task: %w", err)
	}
	return nil
}
