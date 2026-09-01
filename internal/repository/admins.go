// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package repository

import (
	"context"
	"database/sql"
	"strings"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/models"
)

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(
	db *sql.DB,
) *AdminRepository {
	return &AdminRepository{
		db: db,
	}
}

func (r *AdminRepository) Create(
	ctx context.Context,
	email string,
	passwordHash string,
	displayName string,
) (int64, error) {
	email = strings.ToLower(
		strings.TrimSpace(email),
	)

	displayName = strings.TrimSpace(
		displayName,
	)

	result, err := r.db.ExecContext(
		ctx,
		`
			INSERT INTO admins (
				email,
				password_hash,
				display_name
			)
			VALUES (?, ?, ?)
		`,
		email,
		passwordHash,
		displayName,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (r *AdminRepository) GetByID(
	ctx context.Context,
	id int64,
) (models.Admin, error) {
	var admin models.Admin
	var active int

	err := r.db.QueryRowContext(
		ctx,
		`
			SELECT
				id,
				email,
				password_hash,
				display_name,
				is_active,
				created_at,
				updated_at
			FROM admins
			WHERE id = ?
		`,
		id,
	).Scan(
		&admin.ID,
		&admin.Email,
		&admin.PasswordHash,
		&admin.DisplayName,
		&active,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)
	if err != nil {
		return models.Admin{}, err
	}

	admin.IsActive = active == 1

	return admin, nil
}

func (r *AdminRepository) GetByEmail(
	ctx context.Context,
	email string,
) (models.Admin, error) {
	var admin models.Admin
	var active int

	email = strings.ToLower(
		strings.TrimSpace(email),
	)

	err := r.db.QueryRowContext(
		ctx,
		`
			SELECT
				id,
				email,
				password_hash,
				display_name,
				is_active,
				created_at,
				updated_at
			FROM admins
			WHERE email = ?
		`,
		email,
	).Scan(
		&admin.ID,
		&admin.Email,
		&admin.PasswordHash,
		&admin.DisplayName,
		&active,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)
	if err != nil {
		return models.Admin{}, err
	}

	admin.IsActive = active == 1

	return admin, nil
}

func (r *AdminRepository) SetActive(
	ctx context.Context,
	id int64,
	active bool,
) error {
	activeValue := 0

	if active {
		activeValue = 1
	}

	result, err := r.db.ExecContext(
		ctx,
		`
			UPDATE admins
			SET
				is_active = ?,
				updated_at = CURRENT_TIMESTAMP
			WHERE id = ?
		`,
		activeValue,
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
