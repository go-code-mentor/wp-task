package service

import (
	"context"
	"fmt"
)

type Task struct {
	ID          uint32
	Name        string
	Description string
}

type StorageTaskGetter interface {
	Task(id uint32) (Task, error)
}

func GetTask(ctx context.Context, s StorageTaskGetter, id uint32) (Task, error) {
	task, err := s.Task(id)
	if err != nil {
		return task, fmt.Errorf("could not get task: %w", err)
	}
	return task, err
}
