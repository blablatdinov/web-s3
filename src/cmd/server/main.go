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

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/blablatdinov/web-s3/src/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var db *sqlx.DB
var rdb *redis.Client
var ctx = context.Background()

func healthCheck(c *fiber.Ctx) error {
	dbStatus := "ok"
	_, err := db.Exec("SELECT 1")
	if err != nil {
		dbStatus = "fail"
	}
	rdbStatus := "ok"
	err = rdb.Get(ctx, "PING").Err()
	if err.Error() != "redis: nil" {
		rdbStatus = "fail"
	}
	return c.JSON(fiber.Map{
		"app":   "ok",
		"db":    dbStatus,
		"redis": rdbStatus,
	})
}

func databaseDsn() string {
	password := os.Getenv("PG_PASSWORD")
	if password != "" {
		return fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			os.Getenv("PG_HOST"),
			os.Getenv("PG_USER"),
			os.Getenv("PG_PASSWORD"),
			os.Getenv("PG_DBNAME"),
			os.Getenv("PG_PORT"),
		)
	} else {
		return fmt.Sprintf(
			"host=%s user=%s dbname=%s port=%s sslmode=disable",
			os.Getenv("PG_HOST"),
			os.Getenv("PG_USER"),
			os.Getenv("PG_DBNAME"),
			os.Getenv("PG_PORT"),
		)

	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	pgsql, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s dbname=%s sslmode=disable",
		os.Getenv("PG_USERNAME"),
		os.Getenv("PG_DB_NAME"),
	))
	ctx := context.Background()
	if err != nil {
		log.Fatalf("Error connectiing to db: %s\n", err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: "",
		DB:       0,
	})
	region := os.Getenv("S3_REGION")
	sess := session.Must(session.NewSession(&aws.Config{
		Region: &region,
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("S3_ACCESS_KEY"),
			os.Getenv("S3_SECRET_KEY"),
			"",
		),
		Endpoint: aws.String(os.Getenv("S3_ENDPOINT")),
	}))
	s3svc := s3.New(sess)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	app := fiber.New()
	app.Get("/health-check", handlers.HealthCheckCtor(pgsql, rdb, s3svc, ctx).Handle)
	api := app.Group("/api/v1")
	api.Post("/users/auth", handlers.UserAuthCtor(pgsql, os.Getenv("SECRET_KEY")).Handle)
	api.Get("/files", handlers.FilesCtor(pgsql, s3svc).Handle)
	log.Fatal(app.Listen(":8090"))
	fmt.Println("Run server...")
	http.ListenAndServe(":8090", nil)
	err := godotenv.Load()
	if err != nil {
		log.Fatal(
			fmt.Printf("Error loading .env file: %s\n", err),
		)
	}
	db, err = sqlx.Connect("postgres", databaseDsn())
	if err != nil {
		log.Fatal(
			fmt.Printf("Error connecting to postgres: %s\n", err),
		)
	}
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15, // use default DB
	})
	app := fiber.New(fiber.Config{
		Immutable: true,
	})
	app.Get("/health-check", healthCheck)
	app.Listen(":8090")
}
