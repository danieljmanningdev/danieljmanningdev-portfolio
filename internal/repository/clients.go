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
	"fmt"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/models"
)

type ClientRepository struct {
	db *sql.DB
}

func NewClientRepository(db *sql.DB) *ClientRepository {
	return &ClientRepository{
		db: db,
	}
}

func (r *ClientRepository) List(
	ctx context.Context,
) ([]models.Client, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			id,
			name,
			email,
			company,
			status,
			notes,
			created_at,
			updated_at
		FROM clients
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list clients: %w", err)
	}
	defer rows.Close()

	var clients []models.Client

	for rows.Next() {
		var client models.Client

		if err := rows.Scan(
			&client.ID,
			&client.Name,
			&client.Email,
			&client.Company,
			&client.Status,
			&client.Notes,
			&client.CreatedAt,
			&client.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan client: %w", err)
		}

		clients = append(clients, client)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate clients: %w", err)
	}

	return clients, nil
}

func (r *ClientRepository) GetByID(
	ctx context.Context,
	id int64,
) (models.Client, error) {
	var client models.Client

	err := r.db.QueryRowContext(ctx, `
		SELECT
			id,
			name,
			email,
			company,
			status,
			notes,
			created_at,
			updated_at
		FROM clients
		WHERE id = ?
	`, id).Scan(
		&client.ID,
		&client.Name,
		&client.Email,
		&client.Company,
		&client.Status,
		&client.Notes,
		&client.CreatedAt,
		&client.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Client{}, sql.ErrNoRows
		}

		return models.Client{}, fmt.Errorf(
			"get client %d: %w",
			id,
			err,
		)
	}

	return client, nil
}

func (r *ClientRepository) Create(
	ctx context.Context,
	name string,
	email string,
	company string,
	notes string,
) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO clients (
			name,
			email,
			company,
			notes
		)
		VALUES (?, ?, ?, ?)
	`, name, email, company, notes)
	if err != nil {
		return 0, fmt.Errorf("create client: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf(
			"get created client id: %w",
			err,
		)
	}

	return id, nil
}

func (r *ClientRepository) Update(
	ctx context.Context,
	id int64,
	name string,
	email string,
	company string,
	status string,
	notes string,
) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE clients
		SET
			name = ?,
			email = ?,
			company = ?,
			status = ?,
			notes = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, name, email, company, status, notes, id)
	if err != nil {
		return fmt.Errorf(
			"update client %d: %w",
			id,
			err,
		)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"check updated client %d: %w",
			id,
			err,
		)
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *ClientRepository) Delete(
	ctx context.Context,
	id int64,
) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM clients
		WHERE id = ?
	`, id)
	if err != nil {
		return fmt.Errorf(
			"delete client %d: %w",
			id,
			err,
		)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"check deleted client %d: %w",
			id,
			err,
		)
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
