package service

import "fmt"

type Task struct {
	ID          uint32
	Name        string
	Description string
}

type StoreItemReaderI interface {
	GetTask(id uint32) (Task, error)
}

func GetTask(s StoreItemReaderI, id uint32) (Task, error) {
	task, err := s.GetTask(id)
	if err != nil {
		return task, fmt.Errorf("could not get task: %w", err)
	}
	return task, err
}
