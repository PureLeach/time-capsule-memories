package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"time_capsule_memories/internal/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type stubCapsuleRepo struct {
	created *models.CapsuleResponse
	err     error
}

func (s *stubCapsuleRepo) Create(_ context.Context, _ *models.CreateCapsuleRequest) (*models.CapsuleResponse, error) {
	return s.created, s.err
}

func (s *stubCapsuleRepo) ClaimDue(_ context.Context, _ int) ([]*models.CapsuleResponse, error) {
	return nil, nil
}

func (s *stubCapsuleRepo) SetStatus(_ context.Context, _ int, _ string) error { return nil }

func newCapsuleRequest(t *testing.T, body any) *http.Request {
	t.Helper()
	raw, err := json.Marshal(body)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/capsules", bytes.NewReader(raw))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return req
}

func TestCreateCapsule_Created(t *testing.T) {
	e := echo.New()
	req := newCapsuleRequest(t, models.CreateCapsuleRequest{
		SenderName:     "Alice",
		SendAt:         "2099-12-31",
		Message:        "hi",
		RecipientEmail: "alice@example.com",
	})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	want := &models.CapsuleResponse{ID: 7, SenderName: "Alice", Status: "waiting"}
	h := &Handler{capsuleRepo: &stubCapsuleRepo{created: want}}

	require.NoError(t, h.CreateCapsule(c))
	require.Equal(t, http.StatusCreated, rec.Code)
	require.Contains(t, rec.Body.String(), `"id":7`)
}

func TestCreateCapsule_BindFailure(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/capsules", strings.NewReader("not json"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &Handler{capsuleRepo: &stubCapsuleRepo{}}

	require.NoError(t, h.CreateCapsule(c))
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "Invalid request payload")
}

func TestCreateCapsule_RepositoryError(t *testing.T) {
	e := echo.New()
	req := newCapsuleRequest(t, models.CreateCapsuleRequest{
		SenderName:     "Alice",
		SendAt:         "2099-12-31",
		Message:        "hi",
		RecipientEmail: "alice@example.com",
	})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &Handler{capsuleRepo: &stubCapsuleRepo{err: errors.New("db down")}}

	require.NoError(t, h.CreateCapsule(c))
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}
