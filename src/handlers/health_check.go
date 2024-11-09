package handlers

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Handler interface {
	Handle(fiberContext *fiber.Ctx) error
}

type HealthCheck struct {
	pgsql *sqlx.DB
	redis *redis.Client
	s3cfg *s3.Client
	ctx   context.Context
}

func HealthCheckCtor(
	pgsql *sqlx.DB,
	rdb *redis.Client,
	s3cfg *s3.Client,
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
		hndlr.ctx,
		&s3.HeadBucketInput{Bucket: aws.String(os.Getenv("S3_BUCKET"))}, //todo
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
