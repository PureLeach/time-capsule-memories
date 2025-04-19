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
	"strings"
	"time"
	"time_capsule_memories/internal/config"
	"time_capsule_memories/internal/models"
)

// createMessage creates a MIME message with attachments and returns it as a byte slice.
func createMessage(from, subject, body, to string, attachments []models.FileObject) ([]byte, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Email headers
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

	// Adding email body
	part, err := writer.CreatePart(map[string][]string{
		"Content-Type": {"text/plain; charset=UTF-8"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create body part: %v", err)
	}
	part.Write([]byte(body))

	// Adding attachments
	for _, attachment := range attachments {
		// Debugging step: Save attachment as a file locally (can be commented out in production)
		// saveAsFile(attachment)

		// Creating attachment part
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(attachment.FileName)))
		h.Set("Content-Type", attachment.ContentType)
		h.Set("Content-Transfer-Encoding", "base64")
		part, err := writer.CreatePart(h)
		if err != nil {
			return nil, fmt.Errorf("failed to create attachment: %v", err)
		}

		// Base64 encoding for attachment content
		encodedContent := base64.StdEncoding.EncodeToString(attachment.Content)
		_, err = part.Write([]byte(encodedContent))
		if err != nil {
			return nil, fmt.Errorf("failed to write attachment content: %v", err)
		}
	}

	// Finalizing MIME writer
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close writer: %v", err)
	}

	return buf.Bytes(), nil
}

// saveAsFile saves the attachment file locally for debugging purposes.
func saveAsFile(attachment models.FileObject) {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(attachment.FileName), os.ModePerm); err != nil {
		log.Fatalf("Error creating directory %s: %v", filepath.Dir(attachment.FileName), err)
	}
	// Save file locally
	err := os.WriteFile(attachment.FileName, attachment.Content, 0644)
	if err != nil {
		log.Fatalf("Error saving file %s: %v", attachment.FileName, err)
	}
}

// SendEmail sends an email with the provided subject, body, and attachments.
func SendEmail(subject, body, to string, attachments []models.FileObject) error {
	// Get configuration values
	config := config.GetConfig()

	// Create the message (MIME format)
	message, err := createMessage(config.SMTPFrom, subject, body, to, attachments)
	if err != nil {
		return fmt.Errorf("failed to create message: %v", err)
	}

	// Logging attempt to send email
	fmt.Println("Attempting to send email...")

	// Set timeout duration for SMTP operation
	timeout := time.Duration(config.SMTPTimeout) * time.Second

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Channel for sending result
	errCh := make(chan error, 1)

	go func() {
		addr := fmt.Sprintf("%s:%s", config.SMTPHost, config.SMTPPort)

		// Determine if this is MailHog (no TLS/STARTTLS)
		if strings.EqualFold(config.SMTPHost, "mailhog") {
			// Send without TLS and without auth (MailHog accepts open relay on 1025)
			errCh <- smtp.SendMail(addr, nil, config.SMTPFrom, []string{to}, message)
			return
		}

		// For other SMTP hosts, set up auth if credentials provided
		var auth smtp.Auth
		if config.SMTPFrom != "" && config.SMTPPassword != "" {
			auth = smtp.PlainAuth("", config.SMTPFrom, config.SMTPPassword, config.SMTPHost)
		}

		// Send email (will negotiate TLS/STARTTLS if supported)
		errCh <- smtp.SendMail(addr, auth, config.SMTPFrom, []string{to}, message)
	}()

	// Wait for send result or timeout
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
