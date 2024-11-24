package minio_client

import (
	"context"
	"fmt"
	"log"
	"time_capsule_memories/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func MinioClient(config config.Config) (*minio.Client, error) {
	client, err := minio.New(config.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinioAccessKey, config.MinioSecretKey, ""),
		Secure: config.MinioUseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	log.Printf("MinIO client initialized, endpoint: %s", config.MinioEndpoint)
	return client, nil
}

func MinioInit(minioClient *minio.Client, bucketName string) {
	_, err := minioClient.ListBuckets(context.Background())
	if err != nil {
		log.Fatalf("Ошибка подключения к MinIO: %v", err)
	}
	log.Println("Подключение к MinIO успешно установлено")

	exists, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		log.Fatalf("Ошибка при проверке существования бакета %s: %v", bucketName, err)
	}

	if !exists {
		log.Printf("Бакет %s не найден. Пытаемся создать...", bucketName)
		err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: ""})
		if err != nil {
			log.Fatalf("Не удалось создать бакет %s: %v", bucketName, err)
		}
		log.Printf("Бакет %s успешно создан", bucketName)
	} else {
		log.Printf("Бакет %s уже существует", bucketName)
	}
}

// func (m *MinioClient) GeneratePresignedURL(bucket, object string, expiry int64) (string, error) {
// 	presignedURL, err := m.client.PresignedPutObject(nil, bucket, object, expiry, nil)
// 	if err != nil {
// 		return "", err
// 	}
// 	return presignedURL.String(), nil
// }
