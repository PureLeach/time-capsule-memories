package minio_client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"time_capsule_memories/internal/config"
	"time_capsule_memories/internal/models"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Store struct {
	client *minio.Client
	bucket string
}

func New(ctx context.Context, cfg *config.Config) (*Store, error) {
	client, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio client init: %w", err)
	}
	slog.Info("minio client initialized", "endpoint", cfg.MinioEndpoint)

	s := &Store{client: client, bucket: cfg.MinioBucketName}

	if _, err := client.ListBuckets(ctx); err != nil {
		return nil, fmt.Errorf("minio list buckets: %w", err)
	}
	slog.Info("minio connection established")

	exists, err := client.BucketExists(ctx, s.bucket)
	if err != nil {
		return nil, fmt.Errorf("minio bucket lookup %q: %w", s.bucket, err)
	}
	if !exists {
		slog.Info("minio bucket missing; creating", "bucket", s.bucket)
		if err := client.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("minio bucket create %q: %w", s.bucket, err)
		}
		slog.Info("minio bucket created", "bucket", s.bucket)
	} else {
		slog.Info("minio bucket present", "bucket", s.bucket)
	}

	return s, nil
}

func (s *Store) Ping(ctx context.Context) error {
	_, err := s.client.ListBuckets(ctx)
	return err
}

func (s *Store) GeneratePresignedUploadURL(ctx context.Context, objectName string, expiration time.Duration) (string, error) {
	u, err := s.client.PresignedPutObject(ctx, s.bucket, objectName, expiration)
	if err != nil {
		return "", fmt.Errorf("presign put %s: %w", objectName, err)
	}
	return u.String(), nil
}

func (s *Store) GetFilesInDirectory(ctx context.Context, directoryUUID string) ([]models.FileObject, error) {
	objectCh := s.client.ListObjects(ctx, s.bucket, minio.ListObjectsOptions{
		Prefix:    directoryUUID + "/",
		Recursive: false,
	})

	var files []models.FileObject
	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("list object: %w", object.Err)
		}

		obj, err := s.client.GetObject(ctx, s.bucket, object.Key, minio.GetObjectOptions{})
		if err != nil {
			return nil, fmt.Errorf("get object %s: %w", object.Key, err)
		}

		var buffer bytes.Buffer
		if _, err := io.Copy(&buffer, obj); err != nil {
			return nil, fmt.Errorf("read object %s: %w", object.Key, err)
		}

		stat, err := obj.Stat()
		if err != nil {
			return nil, fmt.Errorf("stat object %s: %w", object.Key, err)
		}

		f := strings.Split(object.Key, "/")
		name := f[len(f)-1]
		contentType := stat.ContentType
		if contentType != "" {
			ext := strings.Split(contentType, "/")
			if len(ext) == 2 {
				name = fmt.Sprintf("%s.%s", name, ext[1])
			}
		}

		files = append(files, models.FileObject{
			FileName:    name,
			Content:     buffer.Bytes(),
			ContentType: contentType,
		})
	}

	return files, nil
}
