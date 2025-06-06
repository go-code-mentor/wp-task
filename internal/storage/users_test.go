package storage_test

import (
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func (suite *Suite) TestGetUserLogin() {
	t := suite.T()

	t.Run("success getting login", func(t *testing.T) {
		token := "test-token"
		loginExpected := "test-login"

		data := []struct {
			login string
			token string
		}{
			{"l1", "t1"},
			{loginExpected, token},
			{"l3", "t3"},
			{"l4", "t4"},
		}

		for _, s := range data {
			query := "INSERT INTO users (login) VALUES ($1)"
			res, err := suite.conn.Exec(suite.ctx, query, s.login)
			assert.NoError(t, err)
			assert.Equal(t, int64(1), res.RowsAffected())

			var userID uint64
			query = "SELECT id FROM users WHERE login=$1"
			err = suite.conn.QueryRow(suite.ctx, query, s.login).Scan(&userID)
			assert.NoError(t, err)

			query = "INSERT INTO access_tokens (user_id, token) VALUES ($1, $2)"
			res, err = suite.conn.Exec(suite.ctx, query, userID, s.token)
			assert.NoError(t, err)
			assert.Equal(t, int64(1), res.RowsAffected())
		}

		loginActual, err := suite.storage.GetUserLogin(suite.ctx, token)
		assert.NoError(t, err)
		assert.Equal(t, loginExpected, loginActual)

		_, err = suite.conn.Exec(suite.ctx, "TRUNCATE users RESTART IDENTITY CASCADE")
		assert.NoError(t, err)
	})

	t.Run("getting unexisted login", func(t *testing.T) {
		task, err := suite.storage.GetUserLogin(suite.ctx, "test-token")
		assert.Error(t, pgx.ErrNoRows, err)
		assert.NotNil(t, task)
	})
}
