package storage_test

import (
	"context"
	"errors"
	"github.com/go-code-mentor/wp-task/internal/entities"
	"github.com/go-code-mentor/wp-task/internal/storage"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPgxConn мокает соединение с базой
type MockPgxConn struct {
	mock.Mock
}

func (m *MockPgxConn) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	argsMock := m.Called(ctx, sql, args)
	return argsMock.Get(0).(pgx.Rows), argsMock.Error(1)
}

// MockRows полностью реализует pgx.Rows
type MockRows struct {
	mock.Mock
	currentRow int
	rows       []map[string]interface{}
}

func (m *MockRows) Close() {
	m.Called()
}

func (m *MockRows) Err() error {
	return m.Called().Error(0)
}

func (m *MockRows) CommandTag() pgconn.CommandTag {
	return m.Called().Get(0).(pgconn.CommandTag)
}

func (m *MockRows) Conn() *pgx.Conn {
	return m.Called().Get(0).(*pgx.Conn)
}

func (m *MockRows) FieldDescriptions() []pgconn.FieldDescription {
	return m.Called().Get(0).([]pgconn.FieldDescription)
}

func (m *MockRows) Next() bool {
	m.currentRow++
	return m.currentRow <= len(m.rows)
}

func (m *MockRows) Scan(dest ...interface{}) error {
	row := m.rows[m.currentRow-1]
	for i := 0; i < len(dest); i++ {
		switch d := dest[i].(type) {
		case *int64:
			*d = row["id"].(int64)
		case *string:
			*d = row["name"].(string)
		default:
			return errors.New("unsupported type")
		}
	}
	return nil
}

func (m *MockRows) Values() ([]interface{}, error) {
	row := m.rows[m.currentRow-1]
	return []interface{}{row["id"], row["name"], row["description"]}, nil
}

func (m *MockRows) RawValues() [][]byte {
	row := m.rows[m.currentRow-1]
	return [][]byte{
		[]byte(row["id"].(string)),
		[]byte(row["name"].(string)),
		[]byte(row["description"].(string)),
	}
}

func TestStorage_Tasks(t *testing.T) {
	t.Run("successful retrieval", func(t *testing.T) {
		mockConn := new(MockPgxConn)
		mockRows := &MockRows{
			rows: []map[string]interface{}{
				{"id": int64(1), "name": "Task 1", "description": "Description 1"},
				{"id": int64(2), "name": "Task 2", "description": "Description 2"},
			},
		}

		// Настройка ожиданий для mockRows
		mockRows.On("Close").Return()
		mockRows.On("Err").Return(nil)
		mockRows.On("CommandTag").Return(pgconn.NewCommandTag("SELECT 2"))
		mockRows.On("Conn").Return((*pgx.Conn)(nil))
		mockRows.On("FieldDescriptions").Return([]pgconn.FieldDescription{})
		mockRows.On("RawValues").Return([][]byte{})

		// Настройка ожиданий для mockConn
		mockConn.On("Query",
			mock.Anything,
			"SELECT id, name, description FROM tasks",
			mock.Anything,
		).Return(mockRows, nil)

		s := &storage.Storage{Conn: mockConn}
		tasks, err := s.Tasks(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, []entities.Task{
			{ID: 1, Name: "Task 1", Description: "Description 1"},
			{ID: 2, Name: "Task 2", Description: "Description 2"},
		}, tasks)

		mockConn.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}
