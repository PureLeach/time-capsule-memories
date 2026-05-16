package repository_test

import (
	"context"
	"regexp"
	"testing"
	"time"

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
		AddRow(1, "Alice", now, now, "hi", "alice@example.com", &folder, "in progress").
		AddRow(2, "Bob", now, now, "hello", "bob@example.com", (*string)(nil), "in progress")

	// The 'status = waiting' WHERE matches the partial btree index from
	// migration 20260515120000_index_capsules_dispatch.sql, and the row
	// claim depends on FOR UPDATE SKIP LOCKED for the concurrency guarantee.
	mock.ExpectQuery(`(?s)status = 'waiting'.*FOR UPDATE SKIP LOCKED`).
		WithArgs(100).
		WillReturnRows(rows)

	repo := repository.NewCapsule(mock)
	got, err := repo.ClaimDue(context.Background(), 100)
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.Equal(t, 1, got[0].ID)
	require.Equal(t, "Alice", got[0].SenderName)
	require.Equal(t, "in progress", got[0].Status)
	require.NotNil(t, got[0].FilesFolderUUID)
	require.Equal(t, folder, *got[0].FilesFolderUUID)
	require.Equal(t, 2, got[1].ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCapsule_SetStatus(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE capsules SET status = $1 WHERE id = $2`)).
		WithArgs("done", 42).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	repo := repository.NewCapsule(mock)
	require.NoError(t, repo.SetStatus(context.Background(), 42, "done"))
	require.NoError(t, mock.ExpectationsWereMet())
}
