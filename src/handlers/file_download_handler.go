package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/blablatdinov/web-s3/src/repo"
	"github.com/blablatdinov/web-s3/src/srv"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type FileDownloadHandler struct {
	bucketsRepo repo.BucketsRepo
	userAuthSrv srv.UserAuth
}

func FileDownloadHandlerCtor(bucketsRepo repo.BucketsRepo, userAuthSrv srv.UserAuth) Handler {
	return FileDownloadHandler{
		bucketsRepo: bucketsRepo,
		userAuthSrv: userAuthSrv,
	}
}

func (h FileDownloadHandler) Handle(c *fiber.Ctx) error {
	userID, ok := GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found in context",
		})
	}
	filePath := c.Params("path")
	if filePath == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File path is required",
		})
	}
	filePath = strings.TrimPrefix(filePath, "/")
	decodedPath, err := url.PathUnescape(filePath)
	if err != nil {
		log.Warnf("Failed to decode URL-encoded path '%s': %s", filePath, err)
	} else {
		filePath = decodedPath
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
	fmt.Printf("%s %s\n", bucket.BucketName, filePath)

	result, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket.BucketName),
		Key:    aws.String(filePath),
	})
	if err != nil {
		log.Errorf("Error getting object from S3: %s", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "File not found",
		})
	}
	defer func() {
		if closeErr := result.Body.Close(); closeErr != nil {
			log.Errorf("Error closing S3 object body: %s", closeErr)
		}
	}()
	fileName := filepath.Base(filePath)
	if fileName == "." || fileName == "/" {
		fileName = "file"
	}
	contentType := "application/octet-stream"
	if ext := filepath.Ext(fileName); ext != "" {
		if detectedType := mime.TypeByExtension(ext); detectedType != "" {
			contentType = detectedType
		}
	}
	if result.ContentType != nil && *result.ContentType != "" {
		contentType = *result.ContentType
	}
	c.Set("Content-Type", contentType)
	// Правильно экранируем имя файла для заголовка Content-Disposition
	// Используем RFC 5987 формат для поддержки Unicode символов
	encodedFileName := url.QueryEscape(fileName)
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`, fileName, encodedFileName))
	if result.ContentLength != nil {
		c.Set("Content-Length", fmt.Sprintf("%d", *result.ContentLength))
	}
	_, err = io.Copy(c.Response().BodyWriter(), result.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to stream file",
		})
	}
	return nil
}
