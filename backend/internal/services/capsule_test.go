package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"time_capsule_memories/internal/models"
	"time_capsule_memories/internal/services"

	"github.com/stretchr/testify/require"
)

type statusCall struct {
	id     int
	status string
}

type fakeRepo struct {
	calls        []statusCall
	setStatusErr error
}

func (f *fakeRepo) Create(_ context.Context, _ *models.CreateCapsuleRequest) (*models.CapsuleResponse, error) {
	return nil, nil
}

func (f *fakeRepo) ClaimDue(_ context.Context, _ int) ([]*models.CapsuleResponse, error) {
	return nil, nil
}

func (f *fakeRepo) SetStatus(_ context.Context, id int, status string) error {
	f.calls = append(f.calls, statusCall{id, status})
	return f.setStatusErr
}

type fakeStore struct {
	files []models.FileObject
	err   error
}

func (f *fakeStore) GetFilesInDirectory(_ context.Context, _ string) ([]models.FileObject, error) {
	return f.files, f.err
}

func (f *fakeStore) GeneratePresignedUploadURL(_ context.Context, _ string, _ time.Duration) (string, error) {
	return "", nil
}

func (f *fakeStore) Ping(_ context.Context) error { return nil }

type fakeMailer struct {
	sent int
	err  error
}

func (f *fakeMailer) Send(_ context.Context, _, _, _ string, _ []models.FileObject) error {
	f.sent++
	return f.err
}

func newCapsule() *models.CapsuleResponse {
	folder := "07023417-5079-429d-a113-cbef2ef164d7"
	return &models.CapsuleResponse{
		ID:              42,
		SenderName:      "Alice",
		SendAt:          time.Now(),
		Message:         "hello",
		RecipientEmail:  "alice@example.com",
		FilesFolderUUID: &folder,
		Status:          "in progress",
	}
}

func TestProcess_HappyPath(t *testing.T) {
	repo := &fakeRepo{}
	store := &fakeStore{files: []models.FileObject{{FileName: "a.txt"}}}
	mailer := &fakeMailer{}
	svc := services.NewCapsuleService(repo, store, mailer)

	require.NoError(t, svc.Process(context.Background(), newCapsule()))
	require.Equal(t, 1, mailer.sent)
	require.Equal(t, []statusCall{{42, "done"}}, repo.calls)
}

func TestProcess_MinioError_RevertsToWaiting(t *testing.T) {
	repo := &fakeRepo{}
	store := &fakeStore{err: errors.New("minio boom")}
	mailer := &fakeMailer{}
	svc := services.NewCapsuleService(repo, store, mailer)

	require.Error(t, svc.Process(context.Background(), newCapsule()))
	require.Zero(t, mailer.sent)
	require.Equal(t, []statusCall{{42, "waiting"}}, repo.calls)
}

func TestProcess_SMTPError_RevertsToWaiting(t *testing.T) {
	repo := &fakeRepo{}
	store := &fakeStore{}
	mailer := &fakeMailer{err: errors.New("smtp boom")}
	svc := services.NewCapsuleService(repo, store, mailer)

	require.Error(t, svc.Process(context.Background(), newCapsule()))
	require.Equal(t, 1, mailer.sent)
	require.Equal(t, []statusCall{{42, "waiting"}}, repo.calls)
}

func TestProcess_SetStatusErrorAfterSend_NoRevert(t *testing.T) {
	repo := &fakeRepo{setStatusErr: errors.New("db boom")}
	store := &fakeStore{}
	mailer := &fakeMailer{}
	svc := services.NewCapsuleService(repo, store, mailer)

	require.Error(t, svc.Process(context.Background(), newCapsule()))
	require.Equal(t, 1, mailer.sent)
	require.Equal(t, []statusCall{{42, "done"}}, repo.calls)
}
