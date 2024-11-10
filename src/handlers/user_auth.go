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
	"fmt"
	"strings"
	"time"

	"github.com/blablatdinov/web-s3/src/srv"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/jmoiron/sqlx"
)

type UserAuth struct {
	pgsql     *sqlx.DB
	secretKey string
}

func UserAuthCtor(pgsql *sqlx.DB, secretKey string) Handler {
	return UserAuth{pgsql, secretKey}
}

func (userAuth UserAuth) Handle(fiberContext *fiber.Ctx) error {
	body := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err := fiberContext.BodyParser(&body)
	if err != nil {
		fmt.Printf("Error parsing body. Err=%s\n", err)
	}
	userId := 0
	if err != nil {
		return fiberContext.Status(fiber.StatusInternalServerError).SendString(
			fmt.Sprintf("Error on hash password: '%s'", err),
		)
	}
	user := struct {
		UserId       int    `db:"user_id"`
		PasswordHash string `db:"password_hash"`
	}{}
	err = userAuth.pgsql.Get(
		&user,
		strings.Join([]string{
			"SELECT user_id, password_hash FROM users",
			"WHERE username=$1",
		}, "\n"),
		body.Username,
	)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return fiberContext.Status(fiber.StatusInternalServerError).SendString(
			fmt.Sprintf("Error sql: '%s'", err),
		)
	}
	if user.UserId == 0 {
		return fiberContext.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"details": fmt.Sprintf("User with username '%s' not found", body.Username),
		})
	}
	passValid := srv.CheckPasswordHash(body.Password, user.PasswordHash)
	if !passValid {
		return fiberContext.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"details": "Invalid password",
		})
	}
	claims := jwt.MapClaims{
		"user_id":  userId,
		"username": body.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(userAuth.secretKey))
	if err != nil {
		return fiberContext.Status(fiber.StatusInternalServerError).SendString(
			fmt.Sprintf("Error generate jwt token: '%s'", err),
		)
	}
	return fiberContext.JSON(fiber.Map{
		"access": t,
	})
}
