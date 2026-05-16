package validators_test

import (
	"strings"
	"testing"

	"time_capsule_memories/internal/models"
	"time_capsule_memories/internal/validators"

	"github.com/stretchr/testify/require"
)

func TestValidateStruct(t *testing.T) {
	cases := []struct {
		name    string
		message string
		wantErr string
	}{
		{name: "valid", message: "hello"},
		{name: "empty message", message: "", wantErr: "Message"},
		{name: "oversize message", message: strings.Repeat("x", 4097), wantErr: "Message"},
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
