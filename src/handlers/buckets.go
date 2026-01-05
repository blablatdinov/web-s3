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

	"github.com/blablatdinov/web-s3/src/repo"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type BucketsListHandler struct {
	bucketsRepo repo.BucketsRepo
}

func BucketsListHandlerCtor(bucketsRepo repo.BucketsRepo) Handler {
	return BucketsListHandler{bucketsRepo: bucketsRepo}
}

func (h BucketsListHandler) Handle(c *fiber.Ctx) error {
	userID, ok := GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found in context",
		})
	}
	buckets, err := h.bucketsRepo.List(userID)
	if err != nil {
		if errors.Is(err, repo.ErrSQL) {
			log.Error("Error listing buckets. Err=%s\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error listing buckets",
			})
		}
		log.Error("Unexpected error listing buckets. Err=%s\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	safeBuckets := make([]fiber.Map, 0, len(buckets))
	for _, bucket := range buckets {
		safeBuckets = append(safeBuckets, fiber.Map{
			"bucket_id":     bucket.BucketID,
			"user_id":       bucket.UserID,
			"bucket_name":   bucket.BucketName,
			"access_key_id": bucket.AccessKeyID,
			"region":        bucket.Region,
			"endpoint":      bucket.Endpoint,
			"created_at":    bucket.CreatedAt,
			"updated_at":    bucket.UpdatedAt,
		})
	}
	return c.JSON(fiber.Map{
		"buckets": safeBuckets,
	})
}

type NewBucketHandler struct {
	bucketsRepo repo.BucketsRepo
}

func NewBucketHandlerCtor(bucketsRepo repo.BucketsRepo) Handler {
	return NewBucketHandler{bucketsRepo: bucketsRepo}
}

func (h NewBucketHandler) Handle(c *fiber.Ctx) error {
	userID, ok := GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found in context",
		})
	}
	body := struct {
		BucketName      string  `json:"bucket_name"`
		AccessKeyID     string  `json:"access_key_id"`
		SecretAccessKey string  `json:"secret_access_key"`
		Region          string  `json:"region"`
		Endpoint        *string `json:"endpoint,omitempty"`
	}{}
	err := c.BodyParser(&body)
	if err != nil {
		log.Error("Error parsing body. Err=%s\n", err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	fmt.Printf("body=%v\n", body)
	if body.BucketName == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "bucket_name is required",
		})
	}
	if body.AccessKeyID == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "access_key_id is required",
		})
	}
	if body.SecretAccessKey == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "secret_access_key is required",
		})
	}
	if body.Region == "" {
		body.Region = "us-east-1"
	}
	bucketID, err := h.bucketsRepo.Create(
		userID,
		body.BucketName,
		body.AccessKeyID,
		body.SecretAccessKey,
		body.Region,
		body.Endpoint,
	)
	if err != nil {
		if errors.Is(err, repo.ErrBucketNameAlreadyExists) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": "Bucket name already exists",
			})
		}
		if errors.Is(err, repo.ErrSQL) {
			log.Error("Error creating bucket. Err=%s\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error creating bucket",
			})
		}
		log.Error("Unexpected error creating bucket. Err=%s\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"bucket_id":   bucketID,
		"bucket_name": body.BucketName,
	})
}
