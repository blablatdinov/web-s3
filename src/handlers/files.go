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
	"errors"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/blablatdinov/web-s3/src/repo"
	"github.com/blablatdinov/web-s3/src/srv"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
)

type FilesHandler struct {
	pgsql       *sqlx.DB
	bucketsRepo repo.BucketsRepo
}

func FilesCtor(pgsql *sqlx.DB, bucketsRepo repo.BucketsRepo) Handler {
	return FilesHandler{pgsql: pgsql, bucketsRepo: bucketsRepo}
}

func (h FilesHandler) Handle(c *fiber.Ctx) error {
	userID, ok := GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found in context",
		})
	}
	queries := c.Queries()
	bucketIDStr, exist := queries["bucket_id"]
	if !exist {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bucket_id is required",
		})
	}
	bucketID, err := strconv.Atoi(bucketIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid bucket_id",
		})
	}
	bucket, err := h.bucketsRepo.GetByID(userID, bucketID)
	if err != nil {
		if errors.Is(err, repo.ErrBucketNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Bucket not found",
			})
		}
		log.Error("Error getting bucket. Err=%s\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error getting bucket",
		})
	}
	ctx := context.Background()
	s3Client, err := srv.CreateS3ClientFromBucket(ctx, bucket)
	if err != nil {
		log.Error("Error creating S3 client. Err=%s\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating S3 client",
		})
	}
	path, exist := queries["path"]
	if !exist {
		path = ""
	}
	if path != "" && !strings.HasSuffix(path, "/") {
		path += "/"
	}
	var files []string
	var dirs []string
	resp, err := s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket.BucketName),
		Prefix:    aws.String(path),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		log.Errorf("Failed to list objects: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list objects",
		})
	}

	for _, item := range resp.Contents {
		files = append(files, *item.Key)
	}
	for _, item := range resp.CommonPrefixes {
		dirs = append(dirs, *item.Prefix)
	}
	return c.JSON(fiber.Map{
		"files":       files,
		"directories": dirs,
	})
}
