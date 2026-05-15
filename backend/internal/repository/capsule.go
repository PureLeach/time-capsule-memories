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

// GetCapsulesByToday retrieves all capsules scheduled for today with a "waiting" status.
func GetCapsulesByToday(ctx context.Context) (capsules []*models.CapsuleResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	currentDate := time.Now().Format("2006-01-02")

	query := `
	SELECT id, sender_name, created_at, send_at, message, recipient_email, files_folder_UUID, status
	FROM capsules
	WHERE send_at::date = $1 AND status = 'waiting';
	`

	rows, err := database.DB.Query(
		ctx,
		query,
		currentDate,
	)
	if err != nil {
		slog.Error("failed to query capsules", "error", err)
		return nil, err
	}
	defer rows.Close()

	// Mapping rows to capsule objects
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
			slog.Error("failed to scan capsule row", "error", err)
			return nil, err
		}
		capsules = append(capsules, capsule)
	}

	if err := rows.Err(); err != nil {
		slog.Error("failed to iterate capsule rows", "error", err)
		return nil, err
	}

	return capsules, nil
}

// UpdateCapsuleStatusByID updates the status of a capsule by its ID.
func UpdateCapsuleStatusByID(ctx context.Context, capsuleID int, newStatus string) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	query := `
	UPDATE capsules
	SET status = $1
	WHERE id = $2;
	`

	// Execute the status update
	_, err := database.DB.Exec(
		ctx,
		query,
		newStatus,
		capsuleID,
	)
	if err != nil {
		slog.Error("failed to update capsule status", "capsule_id", capsuleID, "error", err)
		return err
	}

	return nil
}
