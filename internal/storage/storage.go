package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/go-code-mentor/wp-task/internal/entities"
	"github.com/jackc/pgx/v5"
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
	Owner       string `db:"owner"`
}

func (s *Storage) Task(ctx context.Context, id uint64, login string) (entities.Task, error) {
	// Create context with timeout for SQL query
	c, cancel := context.WithTimeout(ctx, rowsRetrieveTimeout)
	defer cancel()

	var taskSQL TaskSQL

	// Run SQL query
	query := `SELECT id, name, description, owner FROM tasks WHERE id=$1 AND owner=$2`
	row, err := s.conn.Query(c, query, id, login)
	if err != nil {
		return entities.Task{}, fmt.Errorf("unable to get query task from storage: %w", err)
	}
	defer row.Close()

	taskSQL, err = pgx.CollectOneRow(row, pgx.RowToStructByPos[TaskSQL])
	if err != nil {
		return entities.Task{}, fmt.Errorf("unable to get parse row to DTO: %w", err)
	}

	// Convert DTO to entity
	task := entities.Task{
		ID:          taskSQL.ID,
		Name:        taskSQL.Name,
		Description: taskSQL.Description,
		Owner:       taskSQL.Owner,
	}
	return task, nil
}

func (s *Storage) Tasks(ctx context.Context, login string) ([]entities.Task, error) {
	// Create context with timeout for SQL query
	c, cancel := context.WithTimeout(ctx, rowsRetrieveTimeout)
	defer cancel()

	// Run SQL query
	query := `SELECT id, name, description, owner FROM tasks WHERE owner=$1`
	rows, err := s.conn.Query(c, query, login)
	if err != nil {
		return nil, fmt.Errorf("unable to get query tasks from storage: %w", err)
	}
	defer rows.Close()

	// Parse SQL query to DTO
	var tasksSQL []TaskSQL
	tasksSQL, err = pgx.CollectRows(rows, pgx.RowToStructByName[TaskSQL])
	if err != nil {
		return nil, fmt.Errorf("unable to get parse rows to DTO: %w", err)
	}

	// Convert DTO to entity
	tasks := make([]entities.Task, len(tasksSQL))
	for i := range tasksSQL {
		tasks[i] = entities.Task{
			ID:          tasksSQL[i].ID,
			Name:        tasksSQL[i].Name,
			Description: tasksSQL[i].Description,
			Owner:       tasksSQL[i].Owner,
		}
	}

	return tasks, nil
}

func (s *Storage) TaskRemove(ctx context.Context, id uint64, login string) error {
	// Create context with timeout for SQL query
	c, cancel := context.WithTimeout(ctx, rowsRetrieveTimeout)
	defer cancel()

	// Run SQL query
	query := `DELETE FROM tasks WHERE id=$1 and owner=$2`
	row, err := s.conn.Exec(c, query, id, login)
	if err != nil {
		return fmt.Errorf("unbale to remove task from storage: %w", err)
	}
	if row.RowsAffected() == 0 {
		return fmt.Errorf("unable to remove task from storage: %w", entities.ErrNoTask)
	}

	return nil
}

func (s *Storage) TaskUpdate(ctx context.Context, task entities.Task, login string) error {
	// Create context with timeout for SQL query
	c, cancel := context.WithTimeout(ctx, rowsRetrieveTimeout)
	defer cancel()

	// Run SQL query
	query := `UPDATE tasks SET name = $1, description = $2 WHERE id = $3 and owner=$4`
	row, err := s.conn.Exec(c, query, task.Name, task.Description, task.ID, login)
	if err != nil {
		return fmt.Errorf("unable to update task in storage: %w", err)
	}
	if row.RowsAffected() == 0 {
		return fmt.Errorf("unable to update task in storage: %w", entities.ErrNoTask)
	}

	return nil
}

func (s *Storage) TaskAdd(ctx context.Context, task entities.Task, login string) (uint64, error) {
	// Create context with timeout for SQL query
	c, cancel := context.WithTimeout(ctx, rowsRetrieveTimeout)
	defer cancel()

	// Convert entity to DTO
	taskSQL := TaskSQL{
		Name:        task.Name,
		Description: task.Description,
		Owner:       login,
	}
	var taskID int64

	// Run SQL query
	query := "INSERT INTO tasks (name, description, owner) VALUES ($1, $2, $3) RETURNING id"
	err := s.conn.QueryRow(c, query, taskSQL.Name, taskSQL.Description, taskSQL.Owner).Scan(&taskID)
	if err != nil {
		return 0, fmt.Errorf("unable to add task to storage: %w", err)
	}

	return uint64(taskID), nil
}
