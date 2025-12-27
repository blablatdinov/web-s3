package repo

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
)

var (
	ErrUsernameAlreadyExist = errors.New("username already exists")
)

type UserSignupRepo interface {
	Create(username, passwordHash string) (int, error)
}

type PgUserSignupRepo struct {
	pgsql *sqlx.DB
}

func PgUserSignupRepoCtor(pgsql *sqlx.DB) UserSignupRepo {
	return PgUserSignupRepo{pgsql}
}

func (u PgUserSignupRepo) Create(username, passwordHash string) (int, error) {
	var userId int
	err := u.pgsql.QueryRow(
		strings.Join([]string{
			"INSERT INTO users (username, password_hash)",
			"VALUES ($1, $2)",
			"RETURNING user_id",
		}, "\n"),
		username, passwordHash,
	).Scan(&userId)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" {
			return 0, ErrUsernameAlreadyExist
		}
		log.Error("Error exec sql query. Err=%s\n", err)
		return 0, ErrSQL
	}
	return userId, nil
}
