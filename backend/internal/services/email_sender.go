package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"path/filepath"
	"strings"
	"time"

	"time_capsule_memories/internal/config"
	"time_capsule_memories/internal/models"
)

type Mailer struct {
	host     string
	port     string
	from     string
	password string
	timeout  time.Duration
}

func NewMailer(cfg *config.Config) *Mailer {
	return &Mailer{
		host:     cfg.SMTPHost,
		port:     cfg.SMTPPort,
		from:     cfg.SMTPFrom,
		password: cfg.SMTPPassword,
		timeout:  time.Duration(cfg.SMTPTimeout) * time.Second,
	}
}

func (m *Mailer) Send(ctx context.Context, subject, body, to string, attachments []models.FileObject) error {
	message, err := buildMessage(m.from, subject, body, to, attachments)
	if err != nil {
		return fmt.Errorf("failed to create message: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		addr := fmt.Sprintf("%s:%s", m.host, m.port)

		if strings.EqualFold(m.host, "mailhog") {
			errCh <- smtp.SendMail(addr, nil, m.from, []string{to}, message)
			return
		}

		var auth smtp.Auth
		if m.from != "" && m.password != "" {
			auth = smtp.PlainAuth("", m.from, m.password, m.host)
		}

		errCh <- smtp.SendMail(addr, auth, m.from, []string{to}, message)
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

func buildMessage(from, subject, body, to string, attachments []models.FileObject) ([]byte, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

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

	part, err := writer.CreatePart(map[string][]string{
		"Content-Type": {"text/plain; charset=UTF-8"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create body part: %v", err)
	}
	part.Write([]byte(body))

	for _, attachment := range attachments {
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(attachment.FileName)))
		h.Set("Content-Type", attachment.ContentType)
		h.Set("Content-Transfer-Encoding", "base64")
		part, err := writer.CreatePart(h)
		if err != nil {
			return nil, fmt.Errorf("failed to create attachment: %v", err)
		}

		encodedContent := base64.StdEncoding.EncodeToString(attachment.Content)
		if _, err := part.Write([]byte(encodedContent)); err != nil {
			return nil, fmt.Errorf("failed to write attachment content: %v", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %v", err)
	}

	return buf.Bytes(), nil
}
