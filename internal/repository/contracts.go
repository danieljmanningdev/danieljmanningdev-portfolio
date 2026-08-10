package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/models"
)

const contractSelectColumns = `
	SELECT
		co.id,
		co.client_id,
		c.name,
		co.project_id,
		COALESCE(p.name, ''),
		co.title,
		co.status,
		co.start_date,
		co.end_date,
		co.value_cents,
		COALESCE(co.notes, ''),
		co.created_at,
		co.updated_at
	FROM contracts AS co
	JOIN clients AS c ON c.id = co.client_id
	LEFT JOIN projects AS p ON p.id = co.project_id
`

type ContractRepository struct {
	db *sql.DB
}

func NewContractRepository(
	db *sql.DB,
) *ContractRepository {
	return &ContractRepository{
		db: db,
	}
}

func (r *ContractRepository) List(
	ctx context.Context,
) ([]models.Contract, error) {
	rows, err := r.db.QueryContext(
		ctx,
		contractSelectColumns+`
			ORDER BY co.created_at DESC, co.id DESC
		`,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"list contracts: %w",
			err,
		)
	}
	defer rows.Close()

	var contracts []models.Contract

	for rows.Next() {
		contract, err := scanContract(rows)
		if err != nil {
			return nil, fmt.Errorf(
				"scan contract: %w",
				err,
			)
		}

		contracts = append(
			contracts,
			contract,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate contracts: %w",
			err,
		)
	}

	return contracts, nil
}

func (r *ContractRepository) ListByClientID(
	ctx context.Context,
	clientID int64,
) ([]models.Contract, error) {
	rows, err := r.db.QueryContext(
		ctx,
		contractSelectColumns+`
			WHERE co.client_id = ?
			ORDER BY co.created_at DESC, co.id DESC
		`,
		clientID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"list contracts for client %d: %w",
			clientID,
			err,
		)
	}
	defer rows.Close()

	var contracts []models.Contract

	for rows.Next() {
		contract, err := scanContract(rows)
		if err != nil {
			return nil, fmt.Errorf(
				"scan contract: %w",
				err,
			)
		}

		contracts = append(
			contracts,
			contract,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate contracts: %w",
			err,
		)
	}

	return contracts, nil
}

func (r *ContractRepository) GetByID(
	ctx context.Context,
	id int64,
) (models.Contract, error) {
	contract, err := scanContract(
		r.db.QueryRowContext(
			ctx,
			contractSelectColumns+`
				WHERE co.id = ?
			`,
			id,
		),
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Contract{}, sql.ErrNoRows
		}

		return models.Contract{}, fmt.Errorf(
			"get contract %d: %w",
			id,
			err,
		)
	}

	return contract, nil
}

func (r *ContractRepository) Create(
	ctx context.Context,
	clientID int64,
	projectID *int64,
	title string,
	status string,
	startDate *time.Time,
	endDate *time.Time,
	valueCents *int64,
	notes string,
) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO contracts (
			client_id,
			project_id,
			title,
			status,
			start_date,
			end_date,
			value_cents,
			notes
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		clientID,
		projectIDValue(projectID),
		title,
		status,
		contractDateValue(startDate),
		contractDateValue(endDate),
		int64PointerValue(valueCents),
		notes,
	)
	if err != nil {
		return 0, fmt.Errorf(
			"create contract: %w",
			err,
		)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf(
			"get created contract id: %w",
			err,
		)
	}

	return id, nil
}

func (r *ContractRepository) Update(
	ctx context.Context,
	id int64,
	clientID int64,
	projectID *int64,
	title string,
	status string,
	startDate *time.Time,
	endDate *time.Time,
	valueCents *int64,
	notes string,
) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE contracts
		SET
			client_id = ?,
			project_id = ?,
			title = ?,
			status = ?,
			start_date = ?,
			end_date = ?,
			value_cents = ?,
			notes = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`,
		clientID,
		projectIDValue(projectID),
		title,
		status,
		contractDateValue(startDate),
		contractDateValue(endDate),
		int64PointerValue(valueCents),
		notes,
		id,
	)
	if err != nil {
		return fmt.Errorf(
			"update contract %d: %w",
			id,
			err,
		)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"check updated contract %d: %w",
			id,
			err,
		)
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *ContractRepository) Cancel(
	ctx context.Context,
	id int64,
) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE contracts
		SET
			status = 'cancelled',
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, id)
	if err != nil {
		return fmt.Errorf(
			"cancel contract %d: %w",
			id,
			err,
		)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"check cancelled contract %d: %w",
			id,
			err,
		)
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

type contractScanner interface {
	Scan(dest ...any) error
}

func scanContract(
	scanner contractScanner,
) (models.Contract, error) {
	var (
		contract  models.Contract
		projectID sql.NullInt64
		startDate nullableContractDate
		endDate   nullableContractDate
		value     sql.NullInt64
	)

	if err := scanner.Scan(
		&contract.ID,
		&contract.ClientID,
		&contract.ClientName,
		&projectID,
		&contract.ProjectName,
		&contract.Title,
		&contract.Status,
		&startDate,
		&endDate,
		&value,
		&contract.Notes,
		&contract.CreatedAt,
		&contract.UpdatedAt,
	); err != nil {
		return models.Contract{}, err
	}

	if projectID.Valid {
		id := projectID.Int64
		contract.ProjectID = &id
	}

	if value.Valid {
		valueCents := value.Int64
		contract.ValueCents = &valueCents
	}

	contract.StartDate = startDate.Time
	contract.EndDate = endDate.Time

	return contract, nil
}

type nullableContractDate struct {
	Time *time.Time
}

func (d *nullableContractDate) Scan(
	value any,
) error {
	switch value := value.(type) {
	case nil:
		d.Time = nil
		return nil

	case time.Time:
		date := time.Date(
			value.Year(),
			value.Month(),
			value.Day(),
			0,
			0,
			0,
			0,
			time.UTC,
		)

		d.Time = &date
		return nil

	case string:
		return d.parse(value)

	case []byte:
		return d.parse(string(value))

	default:
		return fmt.Errorf(
			"scan contract date: unsupported type %T",
			value,
		)
	}
}

func (d *nullableContractDate) parse(
	value string,
) error {
	if value == "" {
		d.Time = nil
		return nil
	}

	formats := []string{
		time.DateOnly,
		time.RFC3339Nano,
		"2006-01-02 15:04:05",
	}

	for _, format := range formats {
		parsed, err := time.Parse(
			format,
			value,
		)
		if err != nil {
			continue
		}

		date := time.Date(
			parsed.Year(),
			parsed.Month(),
			parsed.Day(),
			0,
			0,
			0,
			0,
			time.UTC,
		)

		d.Time = &date
		return nil
	}

	return fmt.Errorf(
		"parse contract date %q",
		value,
	)
}

func contractDateValue(
	value *time.Time,
) any {
	if value == nil {
		return nil
	}

	return value.Format(time.DateOnly)
}

func projectIDValue(
	value *int64,
) any {
	if value == nil {
		return nil
	}

	return *value
}

func int64PointerValue(
	value *int64,
) any {
	if value == nil {
		return nil
	}

	return *value
}
