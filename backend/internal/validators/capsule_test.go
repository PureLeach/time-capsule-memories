package validators_test

import (
	"strings"
	"testing"

	"time_capsule_memories/internal/models"
	"time_capsule_memories/internal/validators"

	"github.com/stretchr/testify/require"
)

func TestValidateCapsule(t *testing.T) {
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
		{name: "empty sender", mutate: func(c *models.CreateCapsuleRequest) { c.SenderName = "" }, wantErr: "SenderName"},
		{name: "bad date format", mutate: func(c *models.CreateCapsuleRequest) { c.SendAt = "31/12/2099" }, wantErr: "YYYY-MM-DD"},
		{name: "past date", mutate: func(c *models.CreateCapsuleRequest) { c.SendAt = "2020-01-01" }, wantErr: "must be in the future"},
		{name: "oversize message", mutate: func(c *models.CreateCapsuleRequest) { c.Message = strings.Repeat("x", 4097) }, wantErr: "cannot exceed 4096"},
		{name: "invalid email", mutate: func(c *models.CreateCapsuleRequest) { c.RecipientEmail = "not-an-email" }, wantErr: "RecipientEmail"},
		{name: "bad uuid", mutate: func(c *models.CreateCapsuleRequest) { c.FilesFolderUUID = &badUUID }, wantErr: "FilesFolderUUID"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := base()
			if tc.mutate != nil {
				tc.mutate(&req)
			}
			err := validators.ValidateCapsule(req)
			if tc.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
