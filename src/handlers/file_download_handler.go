package handlers

import (
	"context"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type FileDownloadHandler struct {
	s3Client *s3.Client
}

func FileDownloadHandlerCtor(s3Client *s3.Client) Handler {
	return FileDownloadHandler{
		s3Client: s3Client,
	}
}

func (h FileDownloadHandler) Handle(fiberContext *fiber.Ctx) error {
	filePath := fiberContext.Params("path")
	if filePath == "" {
		return fiberContext.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File path is required",
		})
	}
	filePath = strings.TrimPrefix(filePath, "/")
	bucket := os.Getenv("S3_BUCKET")
	if bucket == "" {
		return fiberContext.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "S3 bucket is not configured",
		})
	}
	ctx := context.Background()
	result, err := h.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filePath),
	})
	if err != nil {
		log.Errorf("Error getting object from S3: %s", err)
		return fiberContext.Status(fiber.StatusNotFound).JSON(fiber.Map{
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
	fiberContext.Set("Content-Type", contentType)
	fiberContext.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	if result.ContentLength != nil {
		fiberContext.Set("Content-Length", fmt.Sprintf("%d", *result.ContentLength))
	}
	_, err = io.Copy(fiberContext.Response().BodyWriter(), result.Body)
	if err != nil {
		return fiberContext.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to stream file",
		})
	}
	return nil
}
