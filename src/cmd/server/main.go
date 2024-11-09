package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/blablatdinov/web-s3/src/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type CustomEndpointResolver struct {
	URL           string
	SigningRegion string
}

func (e *CustomEndpointResolver) ResolveEndpoint(service, region string, options ...interface{}) (aws.Endpoint, error) {
	if service == s3.ServiceID {
		return aws.Endpoint{
			URL:           e.URL,
			SigningRegion: e.SigningRegion,
		}, nil
	}
	return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
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
	awsCfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(os.Getenv("S3_REGION")),
		config.WithEndpointResolverWithOptions(
			&CustomEndpointResolver{
				URL:           os.Getenv("S3_ENDPOINT"),
				SigningRegion: os.Getenv("S3_REGION"),
			},
		),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv("S3_ACCESS_KEY"),
			os.Getenv("S3_SECRET_KEY"),
			"",
		)),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	app := fiber.New()
	app.Get("/health-check", handlers.HealthCheckCtor(pgsql, rdb, s3.NewFromConfig(awsCfg), ctx).Handle)
	log.Fatal(app.Listen(":8090"))
}
