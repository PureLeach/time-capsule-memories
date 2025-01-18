package minio_client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"reflect"
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

// Получение списка файлов по uuid каталога вместе с содержимым
func GetFilesInDirectory(directoryUUID string) ([]models.FileObject, error) {
	bucketName := config.GetConfig().MinioBucketName
	ctx := context.Background()

	minioClient, err := GetMinioClient()
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении MinIO клиента: %w", err)
	}

	// Создаем канал для листинга объектов
	objectCh := minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    directoryUUID + "/", // Указываем префикс каталога
		Recursive: false,               // Не углубляемся в подкаталоги
	})

	var files []models.FileObject
	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("ошибка при получении объекта из MinIO: %w", object.Err)
		}

		// Получаем объект
		obj, err := minioClient.GetObject(ctx, bucketName, object.Key, minio.GetObjectOptions{})
		if err != nil {
			return nil, fmt.Errorf("ошибка при получении содержимого объекта %s: %w", object.Key, err)
		}

		// Читаем содержимое объекта
		var buffer bytes.Buffer
		if _, err := io.Copy(&buffer, obj); err != nil {
			return nil, fmt.Errorf("ошибка при чтении содержимого объекта %s: %w", object.Key, err)
		}

		// Получаем информацию об объекте
		stat, err := obj.Stat()
		if err != nil {
			return nil, fmt.Errorf("ошибка при получении информации об объекте %s: %w", object.Key, err)
		}

		// Извлекаем имя файла и добавляем тип расширения из Content-Type
		f := strings.Split(object.Key, "/")
		name := f[len(f)-1]
		contentType := stat.ContentType
		fmt.Printf("contentType: %#v Type: %v\n", contentType, reflect.TypeOf(contentType))
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
