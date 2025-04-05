package service

import (
	"context"
	"fmt"
	"github.com/go-code-mentor/wp-task/internal/entities"
)

type StorageTaskAdder interface {
	TaskAdd(ctx context.Context, task entities.Task) error
}

func TaskAdd(ctx context.Context, s StorageTaskAdder, task entities.Task) error {
	if err := s.TaskAdd(ctx, task); err != nil {
		return fmt.Errorf("unable to add task: %w", err)
	}
	return nil
}
