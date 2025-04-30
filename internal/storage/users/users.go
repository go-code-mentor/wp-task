package users

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

const rowsRetrieveTimeout = 10 * time.Second

type UserStorage struct {
	conn *pgx.Conn
}

func New(conn *pgx.Conn) *UserStorage {
	return &UserStorage{
		conn: conn,
	}
}

func (s *UserStorage) GetUserLogin(ctx context.Context, token string) (string, error) {

	c, cancel := context.WithTimeout(ctx, rowsRetrieveTimeout)
	defer cancel()

	var login string

	query := `SELECT login FROM users as u LEFT JOIN access_tokens as t ON u.id = t.user_id WHERE t.token=$1`
	err := s.conn.QueryRow(c, query, token).Scan(&login)
	if err != nil {
		return "", fmt.Errorf("unable to get user by access token: %w", err)
	}

	return login, nil
}
