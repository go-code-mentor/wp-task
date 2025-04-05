package service

import (
	"context"
	"fmt"

	"github.com/go-code-mentor/wp-task/internal/entities"
)

type StorageTaskGetter interface {
	Task(ctx context.Context, id uint64) (entities.Task, error)
}

func Task(ctx context.Context, s StorageTaskGetter, id uint64) (entities.Task, error) {
	task, err := s.Task(ctx, id)
	if err != nil {
		return task, fmt.Errorf("could not get task: %w", err)
	}
	return task, err
}
