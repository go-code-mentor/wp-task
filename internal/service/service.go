package service

import (
	"context"
	"fmt"

	"github.com/go-code-mentor/wp-task/internal/entities"
)

type StorageTaskGetter interface {
	Task(ctx context.Context, id uint64) (entities.Task, error)
}

type StorageTasksGetter interface {
	Tasks(ctx context.Context) ([]entities.Task, error)
}

func Task(ctx context.Context, s StorageTaskGetter, id uint64) (entities.Task, error) {
	task, err := s.Task(ctx, id)
	if err != nil {
		return task, fmt.Errorf("could not get task: %w", err)
	}
	return task, nil
}

func Tasks(ctx context.Context, s StorageTasksGetter) ([]entities.Task, error) {
	tasks, err := s.Tasks(ctx)
	if err != nil {
		return tasks, fmt.Errorf("could not get task: %w", err)
	}
	return tasks, nil
}
