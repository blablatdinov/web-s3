package handlers

import (
	"errors"

	"github.com/blablatdinov/web-s3/src/repo"
	"github.com/blablatdinov/web-s3/src/srv"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type UserSingUpHandler struct {
	srv srv.UserSignupSrv
}

func UserSingUpCtor(srv srv.UserSignupSrv) Handler {
	return UserSingUpHandler{srv}
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
	userId, err := h.srv.Create(body.Username, body.Password)
	if err != nil {
		if errors.Is(err, srv.ErrorHashingPassword) {
			log.Errorf("Error hashing password %s", err.Error())
			return fiberContext.Status(500).JSON(fiber.Map{"details": "Error hashing password"})
		} else if errors.Is(err, repo.ErrUsernameAlreadyExist) {
			return fiberContext.Status(422).JSON(fiber.Map{"details": "Username already exists"})
		} else if errors.Is(err, repo.ErrSQL) {
			log.Error("Error exec sql query. Err=%s\n", err)
			return fiberContext.Status(500).JSON(fiber.Map{"details": "Error exec sql query"})
		}
	}
	return fiberContext.JSON(fiber.Map{
		"user_id":  userId,
		"username": body.Username,
	})
}
