package minio_client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
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
)

// GetMinioClient returns a Singleton instance of the MinIO client.
func GetMinioClient() (*minio.Client, error) {
	// Use sync.Once to ensure the client is initialized only once
	once.Do(func() {
		var err error
		minioClientInstance, err = minio.New(config.GetConfig().MinioEndpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(config.GetConfig().MinioAccessKey, config.GetConfig().MinioSecretKey, ""),
			Secure: config.GetConfig().MinioUseSSL,
		})
		if err != nil {
			log.Fatalf("Error initializing MinIO client: %v", err)
		}

		log.Printf("MinIO client initialized, endpoint: %s", config.GetConfig().MinioEndpoint)
	})

	return minioClientInstance, nil
}

// MinioInit initializes MinIO and creates the bucket if it doesn't exist.
func MinioInit() {
	bucketName := config.GetConfig().MinioBucketName
	minioClient, err := GetMinioClient()
	if err != nil {
		log.Fatalf("Error getting MinIO client: %v", err)
	}

	// Check if MinIO connection is successful
	_, err = minioClient.ListBuckets(context.Background())
	if err != nil {
		log.Fatalf("Error connecting to MinIO: %v", err)
	}
	log.Println("Successfully connected to MinIO")

	// Check if the bucket exists
	exists, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		log.Fatalf("Error checking if bucket %s exists: %v", bucketName, err)
	}

	// Create bucket if it does not exist
	if !exists {
		log.Printf("Bucket %s not found. Attempting to create...", bucketName)
		err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: ""})
		if err != nil {
			log.Fatalf("Failed to create bucket %s: %v", bucketName, err)
		}
		log.Printf("Bucket %s created successfully", bucketName)
	} else {
		log.Printf("Bucket %s already exists", bucketName)
	}
}

// GeneratePresignedUploadURL generates a presigned URL for uploading a file to MinIO.
func GeneratePresignedUploadURL(objectName string, expiration time.Duration) (string, error) {
	bucketName := config.GetConfig().MinioBucketName
	ctx := context.Background()

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
func GetFilesInDirectory(directoryUUID string) ([]models.FileObject, error) {
	bucketName := config.GetConfig().MinioBucketName
	ctx := context.Background()

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
