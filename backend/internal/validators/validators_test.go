package validators_test

import (
	"strings"
	"testing"

	"time_capsule_memories/internal/models"
	"time_capsule_memories/internal/validators"

	"github.com/stretchr/testify/require"
)

func TestValidateStruct_Capsule(t *testing.T) {
	validUUID := "07023417-5079-429d-a113-cbef2ef164d7"
	badUUID := "not-a-uuid"

	base := func() models.CreateCapsuleRequest {
		return models.CreateCapsuleRequest{
			SenderName:      "Alice",
			SendAt:          "2099-12-31",
			Message:         "hello",
			RecipientEmail:  "alice@example.com",
			FilesFolderUUID: &validUUID,
		}
	}

	cases := []struct {
		name    string
		mutate  func(*models.CreateCapsuleRequest)
		wantErr string
	}{
		{name: "valid"},
		{name: "nil files folder is allowed", mutate: func(c *models.CreateCapsuleRequest) { c.FilesFolderUUID = nil }},
		{name: "empty sender", mutate: func(c *models.CreateCapsuleRequest) { c.SenderName = "" }, wantErr: "`sender_name` is required"},
		{name: "bad date format", mutate: func(c *models.CreateCapsuleRequest) { c.SendAt = "31/12/2099" }, wantErr: "YYYY-MM-DD"},
		{name: "past date", mutate: func(c *models.CreateCapsuleRequest) { c.SendAt = "2020-01-01" }, wantErr: "must be a date in the future"},
		{name: "oversize message", mutate: func(c *models.CreateCapsuleRequest) { c.Message = strings.Repeat("x", 4097) }, wantErr: "`message` must be at most 4096"},
		{name: "invalid email", mutate: func(c *models.CreateCapsuleRequest) { c.RecipientEmail = "nope" }, wantErr: "`recipient_email` must be a valid email"},
		{name: "bad uuid", mutate: func(c *models.CreateCapsuleRequest) { c.FilesFolderUUID = &badUUID }, wantErr: "`files_folder_uuid` must be a valid UUID"},
		// Regression: an oversize field other than Message used to be reported valid.
		{name: "oversize sender name", mutate: func(c *models.CreateCapsuleRequest) { c.SenderName = strings.Repeat("x", 101) }, wantErr: "`sender_name` must be at most 100"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := base()
			if tc.mutate != nil {
				tc.mutate(&req)
			}
			err := validators.ValidateStruct(req)
			if tc.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

func TestValidateStruct_ReportsEveryFailingField(t *testing.T) {
	err := validators.ValidateStruct(models.CreateCapsuleRequest{})
	require.Error(t, err)
	for _, field := range []string{"sender_name", "send_at", "message", "recipient_email"} {
		require.Contains(t, err.Error(), field)
	}
}

func TestValidateStruct_Feedback(t *testing.T) {
	cases := []struct {
		name    string
		message string
		wantErr string
	}{
		{name: "valid", message: "hello"},
		{name: "empty message", message: "", wantErr: "`message` is required"},
		{name: "oversize message", message: strings.Repeat("x", 4097), wantErr: "`message` must be at most 4096"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.ValidateStruct(models.CreateFeedbackRequest{Message: tc.message})
			if tc.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

func TestValidateStruct_PresignedURLRequest(t *testing.T) {
	const dir = "07023417-5079-429d-a113-cbef2ef164d7"

	cases := []struct {
		name        string
		directory   string
		contentType string
		wantErr     string
	}{
		{name: "valid", directory: dir, contentType: "image/png"},
		{name: "type is case insensitive", directory: dir, contentType: "IMAGE/JPEG"},
		{name: "empty directory", contentType: "image/png", wantErr: "`directory` is required"},
		{name: "bad uuid", directory: "not-a-uuid", contentType: "image/png", wantErr: "`directory` must be a valid UUID"},
		{name: "missing content type", directory: dir, wantErr: "`content_type` is required"},
		// Browsers execute scripts inside an SVG, so it must never be signed for.
		{name: "svg is rejected", directory: dir, contentType: "image/svg+xml", wantErr: "`content_type` must be one of"},
		{name: "non image is rejected", directory: dir, contentType: "application/pdf", wantErr: "`content_type` must be one of"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.ValidateStruct(models.GeneratePresignedURLRequest{
				Directory:   tc.directory,
				ContentType: tc.contentType,
			})
			if tc.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
