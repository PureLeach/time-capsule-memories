package services

import (
	"fmt"
	"log"
	"time_capsule_memories/internal/models"
	"time_capsule_memories/internal/repository"
)

func ProcessCapsule(capsule *models.CapsuleResponse) error {
	log.Printf("Обработка капсулы с ID: %d", capsule.ID)

	// Проверка наличия значения RecipientTgUsername
	if capsule.RecipientTgUsername != "" {
		log.Printf("Recipient Telegram Username: %s", capsule.RecipientTgUsername)
	}
	if *capsule.FilesFolderUUID != "" {
		log.Printf("Files Folder UUID: %s", *capsule.FilesFolderUUID)
	}

	// Настройка получателя, темы и тела письма
	subject := "Test email without attachments 123"
	body := "This is a test email without any attachments."
	attachments := []string{}
	err := SendEmail(subject, body, capsule.RecipientEmail, attachments)
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	err = repository.UpdateCapsuleStatusByID(capsule.ID, "done")
	if err != nil {
		return fmt.Errorf("error updating status: %v", err)
	}
	log.Println("Обработка капсулы завершена")
	return nil
}
