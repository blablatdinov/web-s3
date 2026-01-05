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

package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
)

var (
	ErrBucketNameAlreadyExists = errors.New("bucket name already exists")
	ErrBucketNotFound          = errors.New("bucket not found")
)

type Bucket struct {
	BucketID        int       `db:"bucket_id"`
	UserID          int       `db:"user_id"`
	BucketName      string    `db:"bucket_name"`
	AccessKeyID     string    `db:"access_key_id"`
	SecretAccessKey string    `db:"secret_access_key"`
	Region          string    `db:"region"`
	Endpoint        *string   `db:"endpoint"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

type BucketsRepo interface {
	List(userID int) ([]Bucket, error)
	GetByID(userID, bucketID int) (*Bucket, error)
	Create(userID int, bucketName, accessKeyID, secretAccessKey, region string, endpoint *string) (int, error)
}

type PgBucketsRepo struct {
	pgsql *sqlx.DB
}

func PgBucketsRepoCtor(pgsql *sqlx.DB) BucketsRepo {
	return PgBucketsRepo{pgsql}
}

func (r PgBucketsRepo) List(userID int) ([]Bucket, error) {
	var buckets []Bucket
	err := r.pgsql.Select(
		&buckets,
		strings.Join([]string{
			"SELECT",
			"  bucket_id,",
			"  user_id,",
			"  bucket_name,",
			"  access_key_id,",
			"  secret_access_key,",
			"  region,",
			"  endpoint,",
			"  created_at,",
			"  updated_at",
			"FROM buckets",
			"WHERE user_id = $1",
			"ORDER BY created_at DESC",
		}, "\n"),
		userID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Bucket{}, nil
		}
		log.Error("Error listing buckets. Err=%s\n", err)
		return nil, fmt.Errorf("%w: %s", ErrSQL, err)
	}
	return buckets, nil
}

func (r PgBucketsRepo) GetByID(userID, bucketID int) (*Bucket, error) {
	var bucket Bucket
	err := r.pgsql.Get(
		&bucket,
		strings.Join([]string{
			"SELECT",
			"  bucket_id,",
			"  user_id,",
			"  bucket_name,",
			"  access_key_id,",
			"  secret_access_key,",
			"  region,",
			"  endpoint,",
			"  created_at,",
			"  updated_at",
			"FROM buckets",
			"WHERE bucket_id = $1 AND user_id = $2",
		}, "\n"),
		bucketID, userID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBucketNotFound
		}
		log.Error("Error getting bucket. Err=%s\n", err)
		return nil, fmt.Errorf("%w: %s", ErrSQL, err)
	}
	return &bucket, nil
}

func (r PgBucketsRepo) Create(userID int, bucketName, accessKeyID, secretAccessKey, region string, endpoint *string) (int, error) {
	var bucketID int
	err := r.pgsql.QueryRow(
		strings.Join([]string{
			"INSERT INTO buckets (user_id, bucket_name, access_key_id, secret_access_key, region, endpoint)",
			"VALUES ($1, $2, $3, $4, $5, $6)",
			"RETURNING bucket_id",
		}, "\n"),
		userID, bucketName, accessKeyID, secretAccessKey, region, endpoint,
	).Scan(&bucketID)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"buckets_user_id_bucket_name_key\"" {
			return 0, ErrBucketNameAlreadyExists
		}
		log.Error("Error creating bucket. Err=%s\n", err)
		return 0, fmt.Errorf("%w: %s", ErrSQL, err)
	}
	return bucketID, nil
}
