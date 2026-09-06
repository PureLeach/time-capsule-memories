package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"time_capsule_memories/internal/models"
	"time_capsule_memories/internal/repository"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestFeedback_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	now := time.Now()
	mock.ExpectQuery(`INSERT INTO users_feedback`).
		WithArgs("great site").
		WillReturnRows(mock.NewRows([]string{"id", "created_at", "message"}).AddRow(3, now, "great site"))

	got, err := repository.NewFeedback(mock).Create(context.Background(), &models.CreateFeedbackRequest{Message: "great site"})
	require.NoError(t, err)
	require.Equal(t, 3, got.ID)
	require.Equal(t, "great site", got.Message)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestFeedback_Create_WrapsDatabaseError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	mock.ExpectQuery(`INSERT INTO users_feedback`).
		WithArgs("hi").
		WillReturnError(errors.New("connection refused"))

	got, err := repository.NewFeedback(mock).Create(context.Background(), &models.CreateFeedbackRequest{Message: "hi"})
	require.Nil(t, got)
	require.ErrorContains(t, err, "insert feedback")
	require.NoError(t, mock.ExpectationsWereMet())
}
