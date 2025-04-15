package service

import (
	"fmt"

	"github.com/go-code-mentor/wp-task/internal/entities"
)

type FakeStorage map[string]entities.Task

func (s *FakeStorage) Task(id string) (entities.Task, error) {
	return entities.Task{}, nil
}

func (s *FakeStorage) Tasks() ([]entities.Task, error) {
	return []entities.Task{}, nil
}

func (s *FakeStorage) TaskRemove(id string) error {
	return nil
}

type TaskStorage interface {
	Task(id string) (entities.Task, error)
	Tasks() ([]entities.Task, error)
	TaskRemove(id string) error
}

func New(storage TaskStorage) *Service {
	return &Service{Storage: storage}
}

type Service struct {
	Storage TaskStorage
}

func (s *Service) Task(id string) (entities.Task, error) {
	task, err := s.Storage.Task(id)
	if err != nil {
		return task, fmt.Errorf("could not get task: %w", err)
	}
	return task, nil
}

func (s *Service) Tasks() ([]entities.Task, error) {
	tasks, err := s.Storage.Tasks()
	if err != nil {
		return tasks, fmt.Errorf("could not get task: %w", err)
	}
	return tasks, nil
}

func (s *Service) TaskRemove(id string) error {
	err := s.Storage.TaskRemove(id)
	if err != nil {
		return fmt.Errorf("could not remove task: %w", err)
	}
	return nil
}
