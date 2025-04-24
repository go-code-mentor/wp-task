package storage

import (
	"context"
	"github.com/go-code-mentor/wp-task/internal/entities"
	"github.com/jackc/pgx/v5"
)

type Storage struct {
	pool *pgx.Conn
}

func New(pool *pgx.Conn) *Storage {
	return &Storage{
		pool: pool,
	}
}

func (s *Storage) Task(ctx context.Context, id uint64) (entities.Task, error) {
	return entities.Task{}, nil
}

func (s *Storage) Tasks(ctx context.Context) ([]entities.Task, error) {
	return []entities.Task{}, nil
}

func (s *Storage) TaskRemove(ctx context.Context, id uint64) error {
	return nil
}

func (s *Storage) TaskUpdate(ctx context.Context, task entities.Task) error {
	return nil
}

func (s *Storage) TaskAdd(ctx context.Context, task entities.Task) error {
	return nil
}
