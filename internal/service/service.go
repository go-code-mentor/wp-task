package service

import (
	"context"
	"fmt"

	tgapi "github.com/go-code-mentor/wp-tg-bot/api"

	"github.com/go-code-mentor/wp-task/internal/entities"
)

type Storage interface {
	TaskStorage
}

type TaskStorage interface {
	Task(ctx context.Context, id uint64, login string) (entities.Task, error)
	Tasks(ctx context.Context, login string) ([]entities.Task, error)
	TaskRemove(ctx context.Context, id uint64, login string) error
	TaskUpdate(ctx context.Context, task entities.Task, login string) error
	TaskAdd(ctx context.Context, task entities.Task, login string) (uint64, error)
}

func New(storage Storage, tgClient tgapi.TgBotClient) *Service {
	return &Service{
		Storage:  storage,
		TgClient: tgClient,
	}
}

type Service struct {
	Storage  Storage
	TgClient tgapi.TgBotClient
}

func (s *Service) Task(ctx context.Context, id uint64, login string) (entities.Task, error) {
	task, err := s.Storage.Task(ctx, id, login)
	if err != nil {
		return task, fmt.Errorf("could not get task: %w", err)
	}
	return task, nil
}

func (s *Service) Tasks(ctx context.Context, login string) ([]entities.Task, error) {
	tasks, err := s.Storage.Tasks(ctx, login)
	if err != nil {
		return tasks, fmt.Errorf("could not get tasks: %w", err)
	}
	return tasks, nil
}

func (s *Service) TaskRemove(ctx context.Context, id uint64, login string) error {
	err := s.Storage.TaskRemove(ctx, id, login)
	if err != nil {
		return fmt.Errorf("could not remove task: %w", err)
	}
	return nil
}

func (s *Service) TaskUpdate(ctx context.Context, task entities.Task, login string) error {
	err := s.Storage.TaskUpdate(ctx, task, login)
	if err != nil {
		return fmt.Errorf("unable to update task: %w", err)
	}
	return nil
}

func (s *Service) TaskAdd(ctx context.Context, task entities.Task, login string) (uint64, error) {
	id, err := s.Storage.TaskAdd(ctx, task, login)
	if err != nil {
		return 0, fmt.Errorf("unable to add task: %w", err)
	}
	return id, nil
}

func (s *Service) TaskSendToTg(ctx context.Context, id uint64, login string) error {
	task, err := s.Task(ctx, id, login)
	if err != nil {
		return fmt.Errorf("sending task to tg error: %w", err)
	}

	_, err = s.TgClient.TaskAdd(context.Background(), &tgapi.TaskAddRequest{
		Id:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		Owner:       task.Owner,
	})

	if err != nil {
		return fmt.Errorf("sending task to tg error: %w", err)
	}

	return nil
}
