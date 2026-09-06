package services

import (
	"encoding/base64"
	"mime"
	"strings"
	"testing"

	"time_capsule_memories/internal/models"

	"github.com/stretchr/testify/require"
)

func TestBuildMessage_EncodesNonASCIISubject(t *testing.T) {
	raw, err := buildMessage("capsule@example.com", "Капсула от Максима", "body", "to@example.com", nil)
	require.NoError(t, err)

	headers, _, ok := strings.Cut(string(raw), "\r\n\r\n")
	require.True(t, ok)

	var subject string
	for _, line := range strings.Split(headers, "\r\n") {
		if after, found := strings.CutPrefix(line, "Subject: "); found {
			subject = after
		}
	}
	require.NotEmpty(t, subject)
	require.NotContains(t, subject, "Капсула", "a raw UTF-8 header renders as mojibake in mail clients")

	decoded, err := new(mime.WordDecoder).DecodeHeader(subject)
	require.NoError(t, err)
	require.Equal(t, "Капсула от Максима", decoded)
}

func TestBuildMessage_IncludesRequiredHeaders(t *testing.T) {
	raw, err := buildMessage("capsule@example.com", "Subject", "body", "to@example.com", nil)
	require.NoError(t, err)

	message := string(raw)
	for _, header := range []string{"From: ", "To: ", "Subject: ", "Date: ", "Message-ID: ", "MIME-Version: 1.0", "Content-Type: multipart/mixed; boundary="} {
		require.Contains(t, message, header)
	}
	require.Contains(t, message, "@example.com>", "Message-ID should be scoped to the sender domain")
}

func TestBuildMessage_WrapsBase64AtRFCLineLength(t *testing.T) {
	content := []byte(strings.Repeat("payload", 600))
	raw, err := buildMessage("capsule@example.com", "s", "b", "to@example.com", []models.FileObject{
		{FileName: "photo.jpg", Content: content, ContentType: "image/jpeg"},
	})
	require.NoError(t, err)

	// The limit applies to encoded body lines; headers may be longer.
	_, body, ok := strings.Cut(string(raw), "\r\n\r\n")
	require.True(t, ok)
	for _, line := range strings.Split(body, "\r\n") {
		if strings.HasPrefix(line, "--") || strings.Contains(line, ": ") {
			continue
		}
		require.LessOrEqual(t, len(line), base64LineLength, "line exceeds the RFC 2045 limit: %q", line)
	}
}

func TestBuildMessage_AttachmentRoundTrips(t *testing.T) {
	content := []byte("binary\x00content")
	raw, err := buildMessage("capsule@example.com", "s", "b", "to@example.com", []models.FileObject{
		{FileName: "photo.jpg", Content: content, ContentType: "image/jpeg"},
	})
	require.NoError(t, err)

	message := string(raw)
	require.Contains(t, message, `filename=photo.jpg`)
	require.Contains(t, message, "Content-Type: image/jpeg")
	require.Contains(t, message, "Content-Transfer-Encoding: base64")

	encoded := base64.StdEncoding.EncodeToString(content)
	require.Contains(t, strings.ReplaceAll(message, "\r\n", ""), encoded)
}

func TestBuildMessage_DefaultsMissingContentType(t *testing.T) {
	raw, err := buildMessage("capsule@example.com", "s", "b", "to@example.com", []models.FileObject{
		{FileName: "unknown.bin", Content: []byte("x")},
	})
	require.NoError(t, err)
	require.Contains(t, string(raw), "Content-Type: application/octet-stream")
}
