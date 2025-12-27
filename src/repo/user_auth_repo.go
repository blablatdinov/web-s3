package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrSQL          = errors.New("sql error")
)

type UserAuthRepo interface {
	UserId(username string) (int, error)
	PasswordHash(username string) (string, error)
}

type PgUserAuthRepo struct {
	pgsql *sqlx.DB
}

func PgUserAuthRepoCtor(pgsql *sqlx.DB) UserAuthRepo {
	return PgUserAuthRepo{pgsql}
}

func (repo PgUserAuthRepo) UserId(username string) (int, error) {
	userId := 0
	err := repo.pgsql.Get(
		&userId,
		strings.Join([]string{
			"SELECT user_id FROM users",
			"WHERE username=$1",
		}, "\n"),
		username,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("%w: %s", ErrUserNotFound, username)
		}
		return 0, fmt.Errorf("%w: %s", ErrSQL, err)
	}
	if userId == 0 {
		return 0, fmt.Errorf("%w: %s", ErrUserNotFound, username)
	}
	return userId, nil
}

func (repo PgUserAuthRepo) PasswordHash(username string) (string, error) {
	passwordHash := ""
	err := repo.pgsql.Get(
		&passwordHash,
		strings.Join([]string{
			"SELECT password_hash FROM users",
			"WHERE username=$1",
		}, "\n"),
		username,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%w: %s", ErrUserNotFound, username)
		}
		return "", fmt.Errorf("%w: %s", ErrSQL, err)
	}
	if passwordHash == "" {
		return "", fmt.Errorf("%w: %s", ErrUserNotFound, username)
	}
	return passwordHash, nil
}
