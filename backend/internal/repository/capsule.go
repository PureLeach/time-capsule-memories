package repository

import (
	"context"
	"log/slog"
	"time"

	"time_capsule_memories/internal/database"
	"time_capsule_memories/internal/models"
)

// dbTimeout caps how long a single repository call may block waiting on the
// database. Derived from the caller's ctx, so request cancellation still wins.
const dbTimeout = 5 * time.Second

// CreateCapsule creates a new capsule in the database and returns the created capsule data.
func CreateCapsule(ctx context.Context, capsule *models.CreateCapsuleRequest) (createdCapsule *models.CapsuleResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	query := `
	INSERT INTO capsules (sender_name, send_at, message, recipient_email, files_folder_UUID)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, sender_name, created_at, send_at, message, recipient_email, files_folder_UUID, status;
    `

	createdCapsule = &models.CapsuleResponse{}

	err = database.DB.QueryRow(
		ctx,
		query,
		capsule.SenderName,
		capsule.SendAt,
		capsule.Message,
		capsule.RecipientEmail,
		capsule.FilesFolderUUID,
	).Scan(
		&createdCapsule.ID,
		&createdCapsule.SenderName,
		&createdCapsule.CreatedAt,
		&createdCapsule.SendAt,
		&createdCapsule.Message,
		&createdCapsule.RecipientEmail,
		&createdCapsule.FilesFolderUUID,
		&createdCapsule.Status,
	)

	if err != nil {
		slog.Error("failed to create capsule", "error", err)
		return nil, err
	}

	return createdCapsule, nil
}

func ClaimDueCapsules(ctx context.Context, limit int) ([]*models.CapsuleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	query := `
	WITH due AS (
	    SELECT id FROM capsules
	    WHERE status = 'waiting' AND send_at < CURRENT_DATE + INTERVAL '1 day'
	    ORDER BY send_at
	    FOR UPDATE SKIP LOCKED
	    LIMIT $1
	)
	UPDATE capsules c
	SET status = 'in progress'
	FROM due
	WHERE c.id = due.id
	RETURNING c.id, c.sender_name, c.created_at, c.send_at, c.message,
	          c.recipient_email, c.files_folder_UUID, c.status;
	`

	rows, err := database.DB.Query(ctx, query, limit)
	if err != nil {
		slog.Error("failed to claim due capsules", "error", err)
		return nil, err
	}
	defer rows.Close()

	var capsules []*models.CapsuleResponse
	for rows.Next() {
		capsule := &models.CapsuleResponse{}
		if err := rows.Scan(
			&capsule.ID,
			&capsule.SenderName,
			&capsule.CreatedAt,
			&capsule.SendAt,
			&capsule.Message,
			&capsule.RecipientEmail,
			&capsule.FilesFolderUUID,
			&capsule.Status,
		); err != nil {
			slog.Error("failed to scan claimed capsule row", "error", err)
			return nil, err
		}
		capsules = append(capsules, capsule)
	}

	if err := rows.Err(); err != nil {
		slog.Error("failed to iterate claimed capsule rows", "error", err)
		return nil, err
	}

	return capsules, nil
}

func SetCapsuleStatus(ctx context.Context, capsuleID int, newStatus string) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	_, err := database.DB.Exec(
		ctx,
		`UPDATE capsules SET status = $1 WHERE id = $2`,
		newStatus,
		capsuleID,
	)
	if err != nil {
		slog.Error("failed to set capsule status",
			"capsule_id", capsuleID,
			"status", newStatus,
			"error", err,
		)
		return err
	}

	return nil
}
