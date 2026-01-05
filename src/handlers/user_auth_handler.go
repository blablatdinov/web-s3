/*
The MIT License (MIT)

Copyright (c) 2024 Almaz Ilaletdinov <a.ilaletdinov@yandex.ru>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE
OR OTHER DEALINGS IN THE SOFTWARE.
*/

package handlers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/blablatdinov/web-s3/src/repo"
	"github.com/blablatdinov/web-s3/src/srv"
	fiber "github.com/gofiber/fiber/v2"
)

type UserAuthHandler struct {
	userAuthSrv srv.UserAuth
}

func UserAuthCtor(u srv.UserAuth) Handler {
	return UserAuthHandler{u}
}

func (userAuth UserAuthHandler) Handle(fiberContext *fiber.Ctx) error {
	body := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err := fiberContext.BodyParser(&body)
	if err != nil {
		fmt.Printf("Error parsing body. Err=%s\n", err)
	}
	t, err := userAuth.userAuthSrv.Jwt(body.Username, body.Password)
	if err != nil {
		return handleAuthError(fiberContext, err)
	}
	return fiberContext.JSON(fiber.Map{
		"access": t,
	})
}

func handleAuthError(fiberContext *fiber.Ctx, err error) error {
	if errors.Is(err, repo.ErrUserNotFound) {
		return fiberContext.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}
	if strings.Contains(err.Error(), "Invalid password") {
		return fiberContext.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}
	if errors.Is(err, repo.ErrSQL) {
		fmt.Printf("Database error: %s\n", err)
		return fiberContext.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	fmt.Printf("Unexpected error: %s\n", err)
	return fiberContext.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "Internal server error",
	})
}
