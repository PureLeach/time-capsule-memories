package minio_client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"time"
	"time_capsule_memories/internal/config"
	"time_capsule_memories/internal/models"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	minioClientInstance *minio.Client
	once                sync.Once
	initErr             error
)

// GetMinioClient returns a Singleton instance of the MinIO client. If the
// underlying constructor fails, the same error is returned on every call so
// callers can decide how to react rather than killing the process.
func GetMinioClient() (*minio.Client, error) {
	once.Do(func() {
		cfg := config.GetConfig()
		minioClientInstance, initErr = minio.New(cfg.MinioEndpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
			Secure: cfg.MinioUseSSL,
		})
		if initErr != nil {
			return
		}
		slog.Info("minio client initialized", "endpoint", cfg.MinioEndpoint)
	})

	return minioClientInstance, initErr
}

// MinioInit verifies connectivity and ensures the configured bucket exists.
// Errors bubble up so main can decide between retry, log-and-exit, or
// degraded-mode operation.
func MinioInit(ctx context.Context) error {
	bucketName := config.GetConfig().MinioBucketName
	minioClient, err := GetMinioClient()
	if err != nil {
		return fmt.Errorf("minio client init: %w", err)
	}

	if _, err := minioClient.ListBuckets(ctx); err != nil {
		return fmt.Errorf("minio list buckets: %w", err)
	}
	slog.Info("minio connection established")

	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("minio bucket lookup %q: %w", bucketName, err)
	}

	if !exists {
		slog.Info("minio bucket missing; creating", "bucket", bucketName)
		if err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("minio bucket create %q: %w", bucketName, err)
		}
		slog.Info("minio bucket created", "bucket", bucketName)
	} else {
		slog.Info("minio bucket present", "bucket", bucketName)
	}
	return nil
}

// GeneratePresignedUploadURL generates a presigned URL for uploading a file to MinIO.
func GeneratePresignedUploadURL(ctx context.Context, objectName string, expiration time.Duration) (string, error) {
	bucketName := config.GetConfig().MinioBucketName

	minioClient, err := GetMinioClient()
	if err != nil {
		return "", fmt.Errorf("error getting MinIO client: %w", err)
	}

	presignedURL, err := minioClient.PresignedPutObject(ctx, bucketName, objectName, expiration)
	if err != nil {
		return "", fmt.Errorf("error generating presigned URL for object %s: %w", objectName, err)
	}

	return presignedURL.String(), nil
}

// GetFilesInDirectory retrieves the list of files in a directory by its UUID, along with the contents.
func GetFilesInDirectory(ctx context.Context, directoryUUID string) ([]models.FileObject, error) {
	bucketName := config.GetConfig().MinioBucketName

	minioClient, err := GetMinioClient()
	if err != nil {
		return nil, fmt.Errorf("error getting MinIO client: %w", err)
	}

	// Create a channel to list objects
	objectCh := minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    directoryUUID + "/", // Specify the directory prefix
		Recursive: false,               // Do not recurse into subdirectories
	})

	var files []models.FileObject
	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("error retrieving object from MinIO: %w", object.Err)
		}

		// Get the object
		obj, err := minioClient.GetObject(ctx, bucketName, object.Key, minio.GetObjectOptions{})
		if err != nil {
			return nil, fmt.Errorf("error getting object %s content: %w", object.Key, err)
		}

		// Read the content of the object
		var buffer bytes.Buffer
		if _, err := io.Copy(&buffer, obj); err != nil {
			return nil, fmt.Errorf("error reading content of object %s: %w", object.Key, err)
		}

		// Get object info
		stat, err := obj.Stat()
		if err != nil {
			return nil, fmt.Errorf("error retrieving information about object %s: %w", object.Key, err)
		}

		// Extract the file name and add the extension from Content-Type
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
