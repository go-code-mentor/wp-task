package service

import (
	"context"
	"fmt"

	"github.com/go-code-mentor/wp-task/internal/entities"
)

type FakeStorage map[uint64]entities.Task

func (s *FakeStorage) Task(ctx context.Context, id uint64) (entities.Task, error) {
	return entities.Task{}, nil
}

func (s *FakeStorage) Tasks(ctx context.Context) ([]entities.Task, error) {
	return []entities.Task{}, nil
}

func (s *FakeStorage) TaskRemove(ctx context.Context, id uint64) error {
	return nil
}

type TaskStorage interface {
	Task(ctx context.Context, id uint64) (entities.Task, error)
	Tasks(ctx context.Context) ([]entities.Task, error)
	TaskRemove(ctx context.Context, id uint64) error
}

func New(storage TaskStorage) *Service {
	return &Service{
		Storage: storage,
	}
}

type Service struct {
	Storage TaskStorage
}

func (s *Service) Task(ctx context.Context, id uint64) (entities.Task, error) {
	task, err := s.Storage.Task(ctx, id)
	if err != nil {
		return task, fmt.Errorf("could not get task: %w", err)
	}
	return task, nil
}

func (s *Service) Tasks(ctx context.Context) ([]entities.Task, error) {
	tasks, err := s.Storage.Tasks(ctx)
	if err != nil {
		return tasks, fmt.Errorf("could not get tasks: %w", err)
	}
	return tasks, nil
}

func (s *Service) TaskRemove(ctx context.Context, id uint64) error {
	err := s.Storage.TaskRemove(ctx, id)
	if err != nil {
		return fmt.Errorf("could not remove task: %w", err)
	}
	return nil
}
