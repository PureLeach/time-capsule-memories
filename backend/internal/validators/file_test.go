package validators_test

import (
	"testing"

	"time_capsule_memories/internal/models"
	"time_capsule_memories/internal/validators"

	"github.com/stretchr/testify/require"
)

func TestValidateGeneratePresignedURLRequest(t *testing.T) {
	cases := []struct {
		name      string
		directory string
		wantErr   string
	}{
		{name: "valid", directory: "07023417-5079-429d-a113-cbef2ef164d7"},
		{name: "empty", directory: "", wantErr: "directory parameter is required"},
		{name: "bad uuid", directory: "not-a-uuid", wantErr: "invalid UUID format"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.ValidateGeneratePresignedURLRequest(&models.GeneratePresignedURLRequest{
				Directory: tc.directory,
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
