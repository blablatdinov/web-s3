package handlers

import (
	"strings"

	"github.com/blablatdinov/web-s3/src/srv"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
)

type UserSingUpHandler struct {
	pgsql *sqlx.DB
}

func UserSingUpCtor(pgsql *sqlx.DB) Handler {
	return UserSingUpHandler{pgsql}
}

func (h UserSingUpHandler) Handle(fiberContext *fiber.Ctx) error {
	body := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err := fiberContext.BodyParser(&body)
	if err != nil {
		log.Error("Error parsing body. Err=%s\n", err)
		return fiberContext.Status(422).JSON(fiber.Map{"details": "Invalid request body"})
	}
	if body.Username == "" {
		return fiberContext.Status(422).JSON(fiber.Map{"details": "Invalid username"})
	}
	hash, err := srv.PswrdCtor(body.Password).Hash()
	if err != nil {
		log.Error("Error hashing password")
		return fiberContext.Status(500).JSON(fiber.Map{"details": "Error hashing password"})
	}
	var userId int
	err = h.pgsql.QueryRow(
		strings.Join([]string{
			"INSERT INTO users (username, password_hash)",
			"VALUES ($1, $2)",
			"RETURNING user_id",
		}, "\n"),
		body.Username, hash,
	).Scan(&userId)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" {
			return fiberContext.Status(422).JSON(fiber.Map{"details": "Username already exists"})
		}
		log.Error("Error exec sql query. Err=%s\n", err)
		return fiberContext.Status(500).JSON(fiber.Map{"details": "Error exec sql query"})
	}
	return fiberContext.JSON(fiber.Map{
		"user_id":  userId,
		"username": body.Username,
	})
}
