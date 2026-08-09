package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/models"
)

const projectSelectColumns = `
	SELECT
		p.id,
		p.client_id,
		c.name,
		p.name,
		COALESCE(p.description, ''),
		p.status,
		p.start_date,
		p.due_date,
		p.created_at,
		p.updated_at
	FROM projects AS p
	JOIN clients AS c ON c.id = p.client_id
`

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{
		db: db,
	}
}

func (r *ProjectRepository) List(
	ctx context.Context,
) ([]models.Project, error) {
	rows, err := r.db.QueryContext(
		ctx,
		projectSelectColumns+`
			ORDER BY p.created_at DESC, p.id DESC
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	defer rows.Close()

	var projects []models.Project

	for rows.Next() {
		project, err := scanProject(rows)
		if err != nil {
			return nil, fmt.Errorf("scan project: %w", err)
		}

		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate projects: %w", err)
	}

	return projects, nil
}

func (r *ProjectRepository) ListByClientID(
	ctx context.Context,
	clientID int64,
) ([]models.Project, error) {
	rows, err := r.db.QueryContext(
		ctx,
		projectSelectColumns+`
			WHERE p.client_id = ?
			ORDER BY p.created_at DESC, p.id DESC
		`,
		clientID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"list projects for client %d: %w",
			clientID,
			err,
		)
	}
	defer rows.Close()

	var projects []models.Project

	for rows.Next() {
		project, err := scanProject(rows)
		if err != nil {
			return nil, fmt.Errorf("scan project: %w", err)
		}

		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate projects: %w", err)
	}

	return projects, nil
}

func (r *ProjectRepository) GetByID(
	ctx context.Context,
	id int64,
) (models.Project, error) {
	project, err := scanProject(
		r.db.QueryRowContext(
			ctx,
			projectSelectColumns+`
				WHERE p.id = ?
			`,
			id,
		),
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Project{}, sql.ErrNoRows
		}

		return models.Project{}, fmt.Errorf(
			"get project %d: %w",
			id,
			err,
		)
	}

	return project, nil
}

func (r *ProjectRepository) Create(
	ctx context.Context,
	clientID int64,
	name string,
	description string,
	status string,
	startDate *time.Time,
	dueDate *time.Time,
) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO projects (
			client_id,
			name,
			description,
			status,
			start_date,
			due_date
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`,
		clientID,
		name,
		description,
		status,
		projectDateValue(startDate),
		projectDateValue(dueDate),
	)
	if err != nil {
		return 0, fmt.Errorf("create project: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf(
			"get created project id: %w",
			err,
		)
	}

	return id, nil
}

func (r *ProjectRepository) Update(
	ctx context.Context,
	id int64,
	clientID int64,
	name string,
	description string,
	status string,
	startDate *time.Time,
	dueDate *time.Time,
) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE projects
		SET
			client_id = ?,
			name = ?,
			description = ?,
			status = ?,
			start_date = ?,
			due_date = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`,
		clientID,
		name,
		description,
		status,
		projectDateValue(startDate),
		projectDateValue(dueDate),
		id,
	)
	if err != nil {
		return fmt.Errorf(
			"update project %d: %w",
			id,
			err,
		)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"check updated project %d: %w",
			id,
			err,
		)
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *ProjectRepository) Archive(
	ctx context.Context,
	id int64,
) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE projects
		SET
			status = 'archived',
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, id)
	if err != nil {
		return fmt.Errorf(
			"archive project %d: %w",
			id,
			err,
		)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"check archived project %d: %w",
			id,
			err,
		)
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

type projectScanner interface {
	Scan(dest ...any) error
}

func scanProject(
	scanner projectScanner,
) (models.Project, error) {
	var (
		project   models.Project
		startDate nullableProjectDate
		dueDate   nullableProjectDate
	)

	if err := scanner.Scan(
		&project.ID,
		&project.ClientID,
		&project.ClientName,
		&project.Name,
		&project.Description,
		&project.Status,
		&startDate,
		&dueDate,
		&project.CreatedAt,
		&project.UpdatedAt,
	); err != nil {
		return models.Project{}, err
	}

	project.StartDate = startDate.Time
	project.DueDate = dueDate.Time

	return project, nil
}

type nullableProjectDate struct {
	Time *time.Time
}

func (d *nullableProjectDate) Scan(value any) error {
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
			"scan project date: unsupported type %T",
			value,
		)
	}
}

func (d *nullableProjectDate) parse(
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
		parsed, err := time.Parse(format, value)
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
		"parse project date %q",
		value,
	)
}

func projectDateValue(value *time.Time) any {
	if value == nil {
		return nil
	}

	return value.Format(time.DateOnly)
}
