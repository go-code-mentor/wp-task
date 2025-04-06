package service

import (
	"context"
	"fmt"

	"github.com/go-code-mentor/wp-task/internal/entities"
)

type StorageTaskGetter interface {
	Task(ctx context.Context, id uint64) (entities.Task, error)
}

type StorageTaskRemover interface {
	TaskRemove(ctx context.Context, id uint64) error
}

func Task(ctx context.Context, s StorageTaskGetter, id uint64) (entities.Task, error) {
	task, err := s.Task(ctx, id)
	if err != nil {
		return task, fmt.Errorf("could not get task: %w", err)
	}
	return task, nil
}

func TaskRemove(ctx context.Context, s StorageTaskRemover, id uint64) error {
	err := s.TaskRemove(ctx, id)
	if err != nil {
		return fmt.Errorf("could not remove task: %w", err)
	}
	return nil
}
