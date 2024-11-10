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
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type Files struct {
	pgsql *sqlx.DB
	s3cfg *s3.S3
}

func FilesCtor(pgsql *sqlx.DB, s3 *s3.S3) Handler {
	return Files{pgsql, s3}
}

func (filesHandler Files) Handle(fiberContext *fiber.Ctx) error {
	var files []string
	var dirs []string
	queries := fiberContext.Queries()
	path, exist := queries["path"]
	if !exist {
		path = ""
	}
	resp, err := filesHandler.s3cfg.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:    aws.String(os.Getenv("S3_BUCKET")),
		Prefix:    aws.String(path),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		log.Fatalf("Failed to list objects: %s", err)
	}
	for _, item := range resp.Contents {
		files = append(files, *item.Key)
	}
	for _, item := range resp.CommonPrefixes {
		dirs = append(dirs, *item.Prefix)
	}
	return fiberContext.JSON(fiber.Map{
		"files":       files,
		"directories": dirs,
	})
}
