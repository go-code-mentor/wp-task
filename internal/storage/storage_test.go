package storage_test

import (
	"context"
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
}

func (suite *Suite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		suite.T().Fatalf("error terminating postgres container: %s", err)
	}
}

func (suite *Suite) TestGetTasks() {
	t := suite.T()

	list, err := suite.storage.Tasks(suite.ctx, "test-user")
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.NotNil(t, list)

}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
