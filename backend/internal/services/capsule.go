package services

import (
	"fmt"
	"log"
	"time_capsule_memories/internal/minio_client"
	"time_capsule_memories/internal/models"
	"time_capsule_memories/internal/repository"
)

func ProcessCapsule(capsule *models.CapsuleResponse) error {
	log.Printf("Обработка капсулы с ID: %d", capsule.ID)

	if *capsule.FilesFolderUUID != "" {
		log.Printf("Files Folder UUID: %s", *capsule.FilesFolderUUID)
	}

	// Получаем файлы из MinIO
	attachments, err := minio_client.GetFilesInDirectory(*capsule.FilesFolderUUID)
	if err != nil {
		log.Fatalf("Ошибка при получении файлов из каталога %s: %v", *capsule.FilesFolderUUID, err)
	}

	// Формируем тему письма и отправляем его по почте вместе с вложениями
	subject := fmt.Sprintf("Вам пришла капсула времени от %s", capsule.SenderName)
	err = SendEmail(subject, capsule.Message, capsule.RecipientEmail, attachments)
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	// Меняем статус капсулы на "done" в БД
	err = repository.UpdateCapsuleStatusByID(capsule.ID, "done")
	if err != nil {
		return fmt.Errorf("error updating status: %v", err)
	}
	log.Println("Обработка капсулы завершена")
	return nil
}
