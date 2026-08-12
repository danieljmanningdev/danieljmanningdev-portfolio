package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/models"
)

type AdminSessionRepository struct {
	db *sql.DB
}

func NewAdminSessionRepository(
	db *sql.DB,
) *AdminSessionRepository {
	return &AdminSessionRepository{
		db: db,
	}
}

func (r *AdminSessionRepository) Create(
	ctx context.Context,
	adminID int64,
	tokenHash string,
	expiresAt time.Time,
) (int64, error) {
	result, err := r.db.ExecContext(
		ctx,
		`
			INSERT INTO admin_sessions (
				admin_id,
				token_hash,
				expires_at
			)
			VALUES (?, ?, ?)
		`,
		adminID,
		tokenHash,
		expiresAt.UTC(),
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (r *AdminSessionRepository) GetByTokenHash(
	ctx context.Context,
	tokenHash string,
) (models.AdminSession, error) {
	var session models.AdminSession

	err := r.db.QueryRowContext(
		ctx,
		`
			SELECT
				id,
				admin_id,
				token_hash,
				expires_at,
				last_seen_at,
				created_at
			FROM admin_sessions
			WHERE token_hash = ?
		`,
		tokenHash,
	).Scan(
		&session.ID,
		&session.AdminID,
		&session.TokenHash,
		&session.ExpiresAt,
		&session.LastSeenAt,
		&session.CreatedAt,
	)
	if err != nil {
		return models.AdminSession{}, err
	}

	return session, nil
}

func (r *AdminSessionRepository) Touch(
	ctx context.Context,
	id int64,
	at time.Time,
) error {
	result, err := r.db.ExecContext(
		ctx,
		`
			UPDATE admin_sessions
			SET last_seen_at = ?
			WHERE id = ?
		`,
		at.UTC(),
		id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *AdminSessionRepository) DeleteByTokenHash(
	ctx context.Context,
	tokenHash string,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`
			DELETE FROM admin_sessions
			WHERE token_hash = ?
		`,
		tokenHash,
	)

	return err
}

func (r *AdminSessionRepository) DeleteByAdminID(
	ctx context.Context,
	adminID int64,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`
			DELETE FROM admin_sessions
			WHERE admin_id = ?
		`,
		adminID,
	)

	return err
}

func (r *AdminSessionRepository) DeleteExpired(
	ctx context.Context,
	now time.Time,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`
			DELETE FROM admin_sessions
			WHERE expires_at <= ?
		`,
		now.UTC(),
	)

	return err
}
