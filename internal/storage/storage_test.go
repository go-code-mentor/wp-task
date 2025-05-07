package storage_test

import (
	"context"
	"github.com/go-code-mentor/wp-task/internal/entities"
	"github.com/go-code-mentor/wp-task/internal/storage"
	"github.com/go-code-mentor/wp-task/internal/testhelper"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type Suite struct {
	suite.Suite
	pgContainer *testhelper.PostgresContainer
	storage     *storage.Storage
	ctx         context.Context
	conn        *pgx.Conn
}

func (suite *Suite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testhelper.CreatePostgresContainer(suite.ctx)
	if err != nil {
		suite.T().Fatalf("%s", err)
	}
	suite.pgContainer = pgContainer

	conn, err := pgx.Connect(context.Background(), suite.pgContainer.ConnectionString)
	if err != nil {
		suite.T().Fatalf("could not connect db: %s", err)
	}

	m, err := migrate.New("file://../../internal/db/migrations", suite.pgContainer.ConnectionString)
	if err != nil {
		suite.T().Fatalf("failed to create migration: %s", err)
	}

	if err := m.Up(); err != nil {
		suite.T().Fatalf("failed to up migration: %s", err)
	}

	repository := storage.New(conn)
	suite.storage = repository

	suite.conn = conn
}

func (suite *Suite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		suite.T().Fatalf("error terminating postgres container: %s", err)
	}
}

func (suite *Suite) TestGetTasks() {
	t := suite.T()

	t.Run("success getting empty tasks list", func(t *testing.T) {
		list, err := suite.storage.Tasks(suite.ctx, "test-user")
		assert.NoError(t, err)
		assert.NotNil(t, list)
	})

	t.Run("success getting tasks list", func(t *testing.T) {
		task1 := entities.Task{
			ID:          1,
			Name:        "test-task-1",
			Description: "test-task-1",
			Owner:       "test-user-1",
		}

		task2 := entities.Task{
			ID:          2,
			Name:        "test-task-2",
			Description: "test-task-2",
			Owner:       "test-user-2",
		}

		query := "INSERT INTO tasks (name, description, owner) VALUES ($1, $2, $3),($4, $5, $6)"
		res, err := suite.conn.Exec(suite.ctx, query, task1.Name, task1.Description, task1.Owner, task2.Name, task2.Description, task2.Owner)
		defer func(conn *pgx.Conn, ctx context.Context, sql string, arguments ...any) {
			_, err := conn.Exec(ctx, sql, arguments)
			assert.NoError(t, err)
		}(suite.conn, suite.ctx, "TRUNCATE tasks")
		assert.NoError(t, err)
		assert.Equal(t, int64(2), res.RowsAffected())

		list1, err := suite.storage.Tasks(suite.ctx, "test-user-1")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(list1))
		assert.Equal(t, task1, list1[0])

		list2, err := suite.storage.Tasks(suite.ctx, "test-user-2")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(list2))
		assert.Equal(t, task2, list2[0])

	})
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
