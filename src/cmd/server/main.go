package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

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
