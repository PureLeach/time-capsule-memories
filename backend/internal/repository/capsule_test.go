package repository_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"time_capsule_memories/internal/models"
	"time_capsule_memories/internal/repository"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestCapsule_ClaimDue(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	folder := "07023417-5079-429d-a113-cbef2ef164d7"
	now := time.Now()
	rows := mock.NewRows([]string{
		"id", "sender_name", "created_at", "send_at", "message",
		"recipient_email", "files_folder_uuid", "status",
	}).
		AddRow(1, "Alice", now, now, "hi", "alice@example.com", &folder, models.StatusInProgress).
		AddRow(2, "Bob", now, now, "hello", "bob@example.com", (*string)(nil), models.StatusInProgress)

	// Both halves are load-bearing: the WHERE matches the partial index, and
	// SKIP LOCKED is the concurrency guarantee.
	mock.ExpectQuery(`(?s)status = 'waiting' AND send_at <= NOW\(\).*FOR UPDATE SKIP LOCKED`).
		WithArgs(100).
		WillReturnRows(rows)

	repo := repository.NewCapsule(mock)
	got, err := repo.ClaimDue(context.Background(), 100)
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.Equal(t, 1, got[0].ID)
	require.Equal(t, "Alice", got[0].SenderName)
	require.Equal(t, models.StatusInProgress, got[0].Status)
	require.NotNil(t, got[0].FilesFolderUUID)
	require.Equal(t, folder, *got[0].FilesFolderUUID)
	require.Equal(t, 2, got[1].ID)
	require.Nil(t, got[1].FilesFolderUUID, "a capsule without attachments must scan back as nil")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCapsule_SetStatus(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE capsules SET status = $1 WHERE id = $2`)).
		WithArgs(models.StatusDone, 42).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	repo := repository.NewCapsule(mock)
	require.NoError(t, repo.SetStatus(context.Background(), 42, models.StatusDone))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCapsule_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	folder := "07023417-5079-429d-a113-cbef2ef164d7"
	now := time.Now()
	req := &models.CreateCapsuleRequest{
		SenderName:      "Alice",
		SendAt:          "2099-12-31",
		Message:         "hi",
		RecipientEmail:  "alice@example.com",
		FilesFolderUUID: &folder,
	}

	mock.ExpectQuery(`INSERT INTO capsules`).
		WithArgs(req.SenderName, req.SendAt, req.Message, req.RecipientEmail, req.FilesFolderUUID).
		WillReturnRows(mock.NewRows([]string{
			"id", "sender_name", "created_at", "send_at", "message",
			"recipient_email", "files_folder_uuid", "status",
		}).AddRow(7, "Alice", now, now, "hi", "alice@example.com", &folder, models.StatusWaiting))

	got, err := repository.NewCapsule(mock).Create(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, 7, got.ID)
	require.Equal(t, models.StatusWaiting, got.Status, "a new capsule starts out waiting for its dispatch date")
	require.NoError(t, mock.ExpectationsWereMet())
}
