package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"os"
	"path/filepath"
	"time"
	"time_capsule_memories/internal/config"
	"time_capsule_memories/internal/models"
)

// Функция для создания MIME-сообщения с вложениями
func createMessage(from, subject, body, to string, attachments []models.FileObject) ([]byte, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Заголовки письма
	headers := map[string]string{
		"From":         from,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": fmt.Sprintf("multipart/mixed; boundary=%s", writer.Boundary()),
	}
	for key, value := range headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	buf.WriteString("\r\n")

	// Добавление тела письма
	part, err := writer.CreatePart(map[string][]string{
		"Content-Type": {"text/plain; charset=UTF-8"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create body part: %v", err)
	}
	part.Write([]byte(body))

	// Добавление вложений
	for _, attachment := range attachments {
		// Сохраняем файл локально для отладки
		// saveAsFile(attachment)

		// Создаем часть для вложения
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(attachment.FileName)))
		h.Set("Content-Type", attachment.ContentType)
		h.Set("Content-Transfer-Encoding", "base64")
		part, err := writer.CreatePart(h)
		if err != nil {
			return nil, fmt.Errorf("failed to create attachment: %v", err)
		}

		// Кодирование вложения в base64
		encodedContent := base64.StdEncoding.EncodeToString(attachment.Content)
		_, err = part.Write([]byte(encodedContent))
		if err != nil {
			return nil, fmt.Errorf("failed to write attachment content: %v", err)
		}
	}

	// Завершаем запись MIME
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close writer: %v", err)
	}

	return buf.Bytes(), nil
}

func saveAsFile(attachment models.FileObject) {
	// Метод для сохранения файла для отладки

	// Создаем директорию, если ее нет
	if err := os.MkdirAll(filepath.Dir(attachment.FileName), os.ModePerm); err != nil {
		log.Fatalf("Ошибка при создании директории %s: %v", filepath.Dir(attachment.FileName), err)
	}
	// Сохраняем файл
	err := os.WriteFile(attachment.FileName, attachment.Content, 0644)
	if err != nil {
		log.Fatalf("Ошибка при сохранении файла %s: %v", attachment.FileName, err)
	}
}

// Функция для отправки письма через SMTP
func SendEmail(subject, body, to string, attachments []models.FileObject) error {
	config := config.GetConfig()
	message, err := createMessage(config.SMTPFrom, subject, body, to, attachments)
	if err != nil {
		return fmt.Errorf("failed to create message: %v", err)
	}

	// Настройка аутентификации
	auth := smtp.PlainAuth("", config.SMTPFrom, config.SMTPPassword, config.SMTPHost)

	// Логирование перед отправкой
	fmt.Println("Attempting to send email...")

	// Создание контекста с таймаутом 7 секунд. Если через 7 секунд операция не завершится, то будет возвращена ошибка
	timeout := time.Duration(config.SMTPTimeout) * time.Second

	// Создание контекста с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Канал для ошибок
	errCh := make(chan error, 1)

	go func() {
		errCh <- smtp.SendMail(config.SMTPHost+":"+config.SMTPPort, auth, config.SMTPFrom, []string{to}, message)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("failed to send email: %v", err)
		}
	case <-ctx.Done():
		return fmt.Errorf("failed to send email: timeout reached")
	}

	return nil
}
