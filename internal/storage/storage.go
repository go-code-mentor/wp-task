package storage

import (
	"context"
	"fmt"
	"github.com/go-code-mentor/wp-task/internal/entities"
	"github.com/jackc/pgx/v5"
	"time"
)

const rowsRetrieveTimeout = 10 * time.Second

type PgxConn interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

type Storage struct {
	Conn PgxConn
}

func New(conn *pgx.Conn) *Storage {
	return &Storage{
		Conn: conn,
	}
}

type TaskSQL struct {
	ID          uint64 `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

func (s *Storage) Task(ctx context.Context, id uint64) (entities.Task, error) {
	return entities.Task{}, nil
}

func (s *Storage) Tasks(ctx context.Context) ([]entities.Task, error) {
	// Create context with timeout for SQL query
	c, cancel := context.WithTimeout(ctx, rowsRetrieveTimeout)
	defer cancel()

	// Run SQL query
	query := `SELECT id, name, description FROM tasks`
	rows, err := s.Conn.Query(c, query)
	if err != nil {
		return nil, fmt.Errorf("unbale to get query tasks from storage: %w", err)
	}

	// Parse SQL query to DTO
	tasksSQL := make([]TaskSQL, 0)
	tasksSQL, err = pgx.CollectRows(rows, pgx.RowToStructByName[TaskSQL])
	if err != nil {
		return nil, fmt.Errorf("unbale to get query tasks from storage: %w", err)
	}

	// Convert DTO to entity
	tasks := make([]entities.Task, len(tasksSQL))
	for i := range tasksSQL {
		tasks[i] = entities.Task{
			ID:          tasksSQL[i].ID,
			Name:        tasksSQL[i].Name,
			Description: tasksSQL[i].Description,
		}
	}

	return tasks, nil
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
