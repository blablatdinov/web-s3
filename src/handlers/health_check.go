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
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
	redis "github.com/redis/go-redis/v9"
)

type HealthCheck struct {
	pgsql *sqlx.DB
	rds   *redis.Client
	s3cfg *s3.S3
	ctx   context.Context
}

func HealthCheckCtor(
	pgsql *sqlx.DB,
	rdb *redis.Client,
	s3cfg *s3.S3,
	ctx context.Context,
) Handler {
	return HealthCheck{pgsql: pgsql, rds: rdb, ctx: ctx, s3cfg: s3cfg}
}

func (hndlr HealthCheck) Handle(fiberContext *fiber.Ctx) error {
	app := true
	postgres := false
	redisAvailable := false
	s3Avaliable := false
	_, err := hndlr.pgsql.Exec("SELECT 1")
	if err == nil {
		postgres = true
	} else {
		log.Warnf("Error on call postgres: \"%s\"", err)
	}
	_, err = hndlr.rds.Ping(hndlr.ctx).Result()
	if err == nil {
		redisAvailable = true
	} else {
		log.Warnf("Error on call redis: \"%s\"", err)
	}
	_, err = hndlr.s3cfg.HeadBucket(
		&s3.HeadBucketInput{Bucket: aws.String(os.Getenv("S3_BUCKET"))},
	)
	if err == nil {
		s3Avaliable = true
	} else {
		log.Warnf("Error on call s3: \"%s\"", err)
	}
	return fiberContext.JSON(fiber.Map{
		"app":      app,
		"postgres": postgres,
		"redis":    redisAvailable,
		"s3":       s3Avaliable,
	})
}
