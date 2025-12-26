package repo

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Files interface {
	List(path string) DTO
}

// TODO: rename
type DTO struct {
	Dirs  []string
	Files []string
}

type S3Files struct {
	s3cfg *s3.Client
}

func S3FilesCtor(s3cfg *s3.Client) Files {
	return S3Files{s3cfg}
}

func (s S3Files) List(path string) DTO {
	var files []string
	var dirs []string
	ctx := context.Background()
	resp, err := s.s3cfg.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:    aws.String(os.Getenv("S3_BUCKET")),
		Prefix:    aws.String(path),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		log.Fatalf("Failed to list objects: %s", err)
	}
	for _, item := range resp.Contents {
		files = append(files, *item.Key)
	}
	for _, item := range resp.CommonPrefixes {
		dirs = append(dirs, *item.Prefix)
	}
	return DTO{
		Dirs:  dirs,
		Files: files,
	}
}
