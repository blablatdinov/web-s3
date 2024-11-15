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
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type HealthCheck struct {
	pgsql *sqlx.DB
	redis *redis.Client
	s3cfg *s3.S3
	ctx   context.Context
}

func HealthCheckCtor(
	pgsql *sqlx.DB,
	rdb *redis.Client,
	s3cfg *s3.S3,
	ctx context.Context,
) Handler {
	return HealthCheck{pgsql: pgsql, redis: rdb, ctx: ctx, s3cfg: s3cfg}
}

func (hndlr HealthCheck) Handle(fiberContext *fiber.Ctx) error {
	app := true
	postgres := false
	redis := false
	s3Avaliable := false
	_, err := hndlr.pgsql.Exec("SELECT 1")
	if err == nil {
		postgres = true
	}
	_, err = hndlr.redis.Ping(hndlr.ctx).Result()
	if err == nil {
		redis = true
	}
	_, err = hndlr.s3cfg.HeadBucket(
		&s3.HeadBucketInput{Bucket: aws.String(os.Getenv("S3_BUCKET"))},
	)
	log.Println(err)
	if err == nil {
		s3Avaliable = true
	}
	return fiberContext.JSON(fiber.Map{
		"app":      app,
		"postgres": postgres,
		"redis":    redis,
		"s3":       s3Avaliable,
	})
}
