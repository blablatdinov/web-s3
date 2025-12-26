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

	"github.com/blablatdinov/web-s3/src/srv"
	fiber "github.com/gofiber/fiber/v2"
)

type UserAuthHandler struct {
	userAuthSrv srv.UserAuthSrv
}

func UserAuthCtor(u srv.UserAuthSrv) Handler {
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
		return fiberContext.Status(fiber.StatusInternalServerError).SendString(
			fmt.Sprintf("Error generate jwt token: '%s'", err),
		)
	}
	return fiberContext.JSON(fiber.Map{
		"access": t,
	})
}
