package minio_client

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
	"time_capsule_memories/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	minioClientInstance *minio.Client
	once                sync.Once
)

// Получаем Singleton экземпляр MinIO клиента
func GetMinioClient() (*minio.Client, error) {
	// Используем sync.Once для гарантии, что клиент создается только один раз
	once.Do(func() {
		var err error
		minioClientInstance, err = minio.New(config.GetConfig().MinioEndpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(config.GetConfig().MinioAccessKey, config.GetConfig().MinioSecretKey, ""),
			Secure: config.GetConfig().MinioUseSSL,
		})
		if err != nil {
			log.Fatalf("Ошибка инициализации MinIO клиента: %v", err)
		}

		log.Printf("MinIO клиент инициализирован, endpoint: %s", config.GetConfig().MinioEndpoint)
	})

	return minioClientInstance, nil
}

// Инициализация MinIO и создание бакета
func MinioInit() {
	bucketName := config.GetConfig().MinioBucketName
	minioClient, err := GetMinioClient()
	if err != nil {
		log.Fatalf("Ошибка при получении MinIO клиента: %v", err)
	}

	_, err = minioClient.ListBuckets(context.Background())
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

// Генерируем presigned URL для загрузки файла в MinIO
func GeneratePresignedUploadURL(objectName string, expiration time.Duration) (string, error) {
	bucketName := config.GetConfig().MinioBucketName
	ctx := context.Background()

	minioClient, err := GetMinioClient()
	if err != nil {
		return "", fmt.Errorf("ошибка при получении MinIO клиента: %w", err)
	}

	presignedURL, err := minioClient.PresignedPutObject(ctx, bucketName, objectName, expiration)
	if err != nil {
		return "", fmt.Errorf("ошибка генерации presigned URL для объекта %s: %w", objectName, err)
	}

	return presignedURL.String(), nil
}
