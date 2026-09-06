package services

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net"
	"net/smtp"
	"net/textproto"
	"path/filepath"
	"strings"
	"time"

	"time_capsule_memories/internal/config"
	"time_capsule_memories/internal/models"
)

// RFC 2045 caps encoded body lines here; some relays reject longer ones.
const base64LineLength = 76

type SMTPMailer struct {
	host     string
	port     string
	from     string
	password string
	timeout  time.Duration
}

func NewSMTPMailer(cfg *config.Config) *SMTPMailer {
	return &SMTPMailer{
		host:     cfg.SMTPHost,
		port:     cfg.SMTPPort,
		from:     cfg.SMTPFrom,
		password: cfg.SMTPPassword,
		timeout:  time.Duration(cfg.SMTPTimeout) * time.Second,
	}
}

func (m *SMTPMailer) Send(ctx context.Context, subject, body, to string, attachments []models.FileObject) error {
	message, err := buildMessage(m.from, subject, body, to, attachments)
	if err != nil {
		return fmt.Errorf("build message: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	// smtp.SendMail takes no context, so it is raced against ctx. The buffer lets
	// the goroutine finish and exit even when the timeout wins.
	errCh := make(chan error, 1)
	go func() {
		var auth smtp.Auth
		if m.password != "" {
			auth = smtp.PlainAuth("", m.from, m.password, m.host)
		}
		errCh <- smtp.SendMail(net.JoinHostPort(m.host, m.port), auth, m.from, []string{to}, message)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("send mail to %s: %w", to, err)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("send mail to %s: %w", to, ctx.Err())
	}
}

// buildMessage assembles a multipart/mixed message. Headers are RFC 2047 encoded
// so non-ASCII names and subjects do not arrive as mojibake.
func buildMessage(from, subject, body, to string, attachments []models.FileObject) ([]byte, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	headers := [][2]string{
		{"From", from},
		{"To", to},
		{"Subject", mime.QEncoding.Encode("utf-8", subject)},
		{"Date", time.Now().Format(time.RFC1123Z)},
		{"Message-ID", messageID(from)},
		{"MIME-Version", "1.0"},
		{"Content-Type", fmt.Sprintf("multipart/mixed; boundary=%s", writer.Boundary())},
	}
	for _, h := range headers {
		fmt.Fprintf(&buf, "%s: %s\r\n", h[0], h[1])
	}
	buf.WriteString("\r\n")

	textPart, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type":              {"text/plain; charset=UTF-8"},
		"Content-Transfer-Encoding": {"base64"},
	})
	if err != nil {
		return nil, fmt.Errorf("create body part: %w", err)
	}
	if err := writeBase64(textPart, []byte(body)); err != nil {
		return nil, fmt.Errorf("write body: %w", err)
	}

	for _, attachment := range attachments {
		name := filepath.Base(attachment.FileName)
		contentType := attachment.ContentType
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		part, err := writer.CreatePart(textproto.MIMEHeader{
			"Content-Type":              {contentType},
			"Content-Transfer-Encoding": {"base64"},
			"Content-Disposition":       {mime.FormatMediaType("attachment", map[string]string{"filename": name})},
		})
		if err != nil {
			return nil, fmt.Errorf("create attachment part %q: %w", name, err)
		}
		if err := writeBase64(part, attachment.Content); err != nil {
			return nil, fmt.Errorf("write attachment %q: %w", name, err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close multipart writer: %w", err)
	}

	return buf.Bytes(), nil
}

func writeBase64(w io.Writer, content []byte) error {
	encoded := base64.StdEncoding.EncodeToString(content)
	for len(encoded) > base64LineLength {
		if _, err := io.WriteString(w, encoded[:base64LineLength]+"\r\n"); err != nil {
			return err
		}
		encoded = encoded[base64LineLength:]
	}
	_, err := io.WriteString(w, encoded)
	return err
}

func messageID(from string) string {
	domain := "localhost"
	if _, host, ok := strings.Cut(from, "@"); ok && host != "" {
		domain = host
	}

	var random [16]byte
	if _, err := rand.Read(random[:]); err != nil {
		// crypto/rand only fails catastrophically; a timestamp still works here.
		return fmt.Sprintf("<%d@%s>", time.Now().UnixNano(), domain)
	}
	return fmt.Sprintf("<%s@%s>", hex.EncodeToString(random[:]), domain)
}
