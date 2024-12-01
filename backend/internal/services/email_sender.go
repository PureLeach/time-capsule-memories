package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"time"
	"time_capsule_memories/internal/config"
)

// Функция для создания MIME-сообщения с вложениями
func createMessage(from, subject, body, to string, attachments []string) ([]byte, error) {
	var buf strings.Builder
	writer := multipart.NewWriter(&buf)

	// Запись заголовков письма
	headers := map[string]string{
		"From":    from,
		"To":      to,
		"Subject": subject,
	}

	// Запись заголовков
	for key, value := range headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}

	// Тип контента - multipart
	buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", writer.Boundary()))
	buf.WriteString("\r\n")

	// Тело письма
	part, err := writer.CreatePart(map[string][]string{
		"Content-Type": {"text/plain; charset=UTF-8"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create body part: %v", err)
	}
	part.Write([]byte(body))

	// Вложения
	for _, file := range attachments {
		err := attachFile(writer, file)
		if err != nil {
			return nil, fmt.Errorf("failed to attach file %s: %v", file, err)
		}
	}

	// Закрытие MIME writer
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close writer: %v", err)
	}

	// Возвращаем сформированное письмо
	return []byte(buf.String()), nil
}

// Функция для прикрепления файла
func attachFile(writer *multipart.Writer, filename string) error {
	// Открытие файла с использованием os.ReadFile
	file, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", filename, err)
	}

	// Создание части для файла
	part, err := writer.CreateFormFile("attachment", filepath.Base(filename))
	if err != nil {
		return fmt.Errorf("failed to create form file: %v", err)
	}

	// Запись содержимого файла в часть
	_, err = part.Write(file)
	if err != nil {
		return fmt.Errorf("failed to write file data: %v", err)
	}

	return nil
}

// Функция для отправки письма через SMTP с таймаутом
func SendEmail(subject, body, to string, attachments []string) error {
	config := config.GetConfig()
	// Создание MIME-сообщения
	message, err := createMessage(config.SMTPFrom, subject, body, to, attachments)
	if err != nil {
		return fmt.Errorf("failed to create message: %v", err)
	}

	// Настройка аутентификации
	auth := smtp.PlainAuth("", config.SMTPFrom, config.SMTPPassword, config.SMTPHost)

	// Логирование перед отправкой
	fmt.Println("Attempting to send email...")

	// Создание контекста с таймаутом 7 секунд
	timeout := time.Duration(config.SMTPTimeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Канал для отправки письма
	errCh := make(chan error, 1)

	// Отправка письма в отдельной горутине
	go func() {
		errCh <- smtp.SendMail(config.SMTPHost+":"+config.SMTPPort, auth, config.SMTPFrom, []string{to}, message)
	}()

	// Ожидание результата или истечения таймаута
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
