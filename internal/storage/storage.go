package storage

import (
	"context"
	"fmt"
	"github.com/go-code-mentor/wp-task/internal/entities"
	"github.com/jackc/pgx/v5"
	"time"
)

const rowsRetrieveTimeout = 10 * time.Second

type Storage struct {
	conn *pgx.Conn
}

func New(conn *pgx.Conn) *Storage {
	return &Storage{
		conn: conn,
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
	rows, err := s.conn.Query(c, query)
	if err != nil {
		return nil, fmt.Errorf("unbale to get query tasks from storage: %w", err)
	}
	defer rows.Close()

	// Parse SQL query to DTO
	var tasksSQL []TaskSQL
	tasksSQL, err = pgx.CollectRows(rows, pgx.RowToStructByName[TaskSQL])
	if err != nil {
		return nil, fmt.Errorf("unbale to get parse rows to DTO: %w", err)
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
	// Create context with timeout for SQL query
	c, cancel := context.WithTimeout(ctx, rowsRetrieveTimeout)
	defer cancel()

	// Convert entity to DTO
	taskSQL := TaskSQL{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description,
	}

	// Run SQL query
	query := "INSERT INTO tasks (name, description, created_at, updated_at) VALUES ($1, $2, now(), now())"
	_, err := s.conn.Exec(c, query, taskSQL.Name, taskSQL.Description)
	if err != nil {
		return fmt.Errorf("unable to add task to storage: %w", err)
	}

	return nil
}
