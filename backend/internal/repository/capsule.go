package repository

import (
	"context"
	"fmt"

	"time_capsule_memories/internal/models"
)

type Capsule struct {
	pool dbPool
}

func NewCapsule(pool dbPool) *Capsule {
	return &Capsule{pool: pool}
}

func (r *Capsule) Create(ctx context.Context, capsule *models.CreateCapsuleRequest) (*models.CapsuleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	const query = `
	INSERT INTO capsules (sender_name, send_at, message, recipient_email, files_folder_UUID)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, sender_name, created_at, send_at, message, recipient_email, files_folder_UUID, status;
	`

	created := &models.CapsuleResponse{}
	err := r.pool.QueryRow(
		ctx,
		query,
		capsule.SenderName,
		capsule.SendAt,
		capsule.Message,
		capsule.RecipientEmail,
		capsule.FilesFolderUUID,
	).Scan(
		&created.ID,
		&created.SenderName,
		&created.CreatedAt,
		&created.SendAt,
		&created.Message,
		&created.RecipientEmail,
		&created.FilesFolderUUID,
		&created.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("insert capsule: %w", err)
	}

	return created, nil
}

// ClaimDue atomically moves up to limit due capsules to 'in progress'. SKIP
// LOCKED lets several dispatchers run without both claiming the same row.
func (r *Capsule) ClaimDue(ctx context.Context, limit int) ([]*models.CapsuleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	const query = `
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

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("claim due capsules: %w", err)
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
			return nil, fmt.Errorf("scan claimed capsule: %w", err)
		}
		capsules = append(capsules, capsule)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate claimed capsules: %w", err)
	}

	return capsules, nil
}

func (r *Capsule) SetStatus(ctx context.Context, capsuleID int, status models.CapsuleStatus) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	if _, err := r.pool.Exec(ctx, `UPDATE capsules SET status = $1 WHERE id = $2`, status, capsuleID); err != nil {
		return fmt.Errorf("set capsule %d status to %q: %w", capsuleID, status, err)
	}

	return nil
}
