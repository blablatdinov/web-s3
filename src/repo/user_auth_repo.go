package repo

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
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
	if err != nil && err.Error() != "sql: no rows in result set" {
		return 0, errors.New(fmt.Sprintf("Error sql: '%s'", err))
	}
	if userId == 0 {
		return 0, errors.New(fmt.Sprintf("User with username '%s' not found", username))
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
	if err != nil && err.Error() != "sql: no rows in result set" {
		return "", errors.New(fmt.Sprintf("Error sql: '%s'", err))
	}
	if passwordHash == "" {
		return "", errors.New(fmt.Sprintf("User with username '%s' not found", username))
	}
	return passwordHash, nil
}
