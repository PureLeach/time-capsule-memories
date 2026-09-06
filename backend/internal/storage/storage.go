// Package storage wraps the S3-compatible object store used for capsule attachments.
package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"path"
	"strings"
	"time"

	"time_capsule_memories/internal/config"
	"time_capsule_memories/internal/models"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Upload limits, enforced by the storage server through the signed POST policy
// rather than by the client.
const (
	MaxAttachmentBytes  = 5 << 20
	imageTypePrefix     = "image/"
	maxAttachmentsRead  = 16
	maxTotalAttachBytes = 32 << 20
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
		return nil, fmt.Errorf("object store client init: %w", err)
	}

	s := &Store{client: client, bucket: cfg.MinioBucketName}

	exists, err := client.BucketExists(ctx, s.bucket)
	if err != nil {
		return nil, fmt.Errorf("object store bucket lookup %q: %w", s.bucket, err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("object store bucket create %q: %w", s.bucket, err)
		}
		slog.Info("object store bucket created", "bucket", s.bucket)
	}

	slog.Info("object store ready", "endpoint", cfg.MinioEndpoint, "bucket", s.bucket)
	return s, nil
}

func (s *Store) Ping(ctx context.Context) error {
	_, err := s.client.BucketExists(ctx, s.bucket)
	return err
}

// PresignUpload returns a multipart/form-data target whose signed policy pins the
// key, the content type and a MaxAttachmentBytes ceiling. Fields must be sent
// back verbatim.
func (s *Store) PresignUpload(ctx context.Context, directoryUUID, contentType string, expiration time.Duration) (*models.PresignedUpload, error) {
	policy := minio.NewPostPolicy()

	if err := policy.SetBucket(s.bucket); err != nil {
		return nil, fmt.Errorf("post policy bucket: %w", err)
	}
	if err := policy.SetKey(path.Join(directoryUUID, uuid.NewString())); err != nil {
		return nil, fmt.Errorf("post policy key: %w", err)
	}
	if err := policy.SetExpires(time.Now().UTC().Add(expiration)); err != nil {
		return nil, fmt.Errorf("post policy expiry: %w", err)
	}
	// Pinned exactly, not by prefix, so the stored object carries a usable type —
	// it becomes the attachment's type in the delivered email.
	if err := policy.SetContentType(contentType); err != nil {
		return nil, fmt.Errorf("post policy content type: %w", err)
	}
	if err := policy.SetContentLengthRange(1, MaxAttachmentBytes); err != nil {
		return nil, fmt.Errorf("post policy size range: %w", err)
	}

	u, formData, err := s.client.PresignedPostPolicy(ctx, policy)
	if err != nil {
		return nil, fmt.Errorf("presign post policy: %w", err)
	}

	return &models.PresignedUpload{URL: u.String(), Fields: formData}, nil
}

// GetFilesInDirectory reads the directory's objects into memory. The caller
// base64-encodes them into one message, so the read is capped to keep an
// oversized folder from exhausting the heap.
func (s *Store) GetFilesInDirectory(ctx context.Context, directoryUUID string) ([]models.FileObject, error) {
	objects := s.client.ListObjects(ctx, s.bucket, minio.ListObjectsOptions{
		Prefix:    directoryUUID + "/",
		Recursive: false,
	})

	var (
		files      []models.FileObject
		totalBytes int64
	)
	for object := range objects {
		if object.Err != nil {
			return nil, fmt.Errorf("list objects in %s: %w", directoryUUID, object.Err)
		}
		if len(files) >= maxAttachmentsRead {
			slog.Warn("attachment count limit reached; ignoring the rest",
				"folder_uuid", directoryUUID, "limit", maxAttachmentsRead)
			break
		}
		if totalBytes+object.Size > maxTotalAttachBytes {
			slog.Warn("attachment size budget reached; ignoring the rest",
				"folder_uuid", directoryUUID, "limit_bytes", maxTotalAttachBytes)
			break
		}

		file, err := s.readObject(ctx, object.Key)
		if err != nil {
			return nil, err
		}
		totalBytes += int64(len(file.Content))
		files = append(files, *file)
	}

	return files, nil
}

func (s *Store) readObject(ctx context.Context, key string) (*models.FileObject, error) {
	obj, err := s.client.GetObject(ctx, s.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("get object %s: %w", key, err)
	}
	defer func() { _ = obj.Close() }()

	stat, err := obj.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat object %s: %w", key, err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, io.LimitReader(obj, MaxAttachmentBytes)); err != nil {
		return nil, fmt.Errorf("read object %s: %w", key, err)
	}

	return &models.FileObject{
		FileName:    attachmentName(key, stat.ContentType),
		Content:     buf.Bytes(),
		ContentType: stat.ContentType,
	}, nil
}

// attachmentName gives the opaque object key an extension, so mail clients show
// it as an image rather than an unnamed blob.
func attachmentName(key, contentType string) string {
	name := path.Base(key)
	subtype, ok := strings.CutPrefix(contentType, imageTypePrefix)
	if !ok || subtype == "" {
		return name
	}
	if i := strings.IndexAny(subtype, "; "); i >= 0 {
		subtype = subtype[:i]
	}
	return name + "." + subtype
}
