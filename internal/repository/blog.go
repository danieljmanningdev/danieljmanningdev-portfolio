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

type BlogRepository struct {
	db *sql.DB
}

func NewBlogRepository(db *sql.DB) *BlogRepository {
	return &BlogRepository{
		db: db,
	}
}

func (r *BlogRepository) ListPublished(
	ctx context.Context,
) ([]models.BlogPost, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			id,
			title,
			slug,
			excerpt,
			status,
			published_at,
			created_at,
			updated_at
		FROM blog_posts
		WHERE status = 'published'
		ORDER BY published_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("list published blog posts: %w", err)
	}
	defer rows.Close()

	var posts []models.BlogPost

	for rows.Next() {
		var post models.BlogPost

		if err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Slug,
			&post.Excerpt,
			&post.Status,
			&post.PublishedAt,
			&post.CreatedAt,
			&post.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan blog post: %w", err)
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate blog posts: %w", err)
	}

	return posts, nil
}

func (r *BlogRepository) GetBySlug(
	ctx context.Context,
	slug string,
) (models.BlogPost, error) {
	var post models.BlogPost

	err := r.db.QueryRowContext(ctx, `
		SELECT
			id,
			title,
			slug,
			excerpt,
			content,
			status,
			published_at,
			created_at,
			updated_at
		FROM blog_posts
		WHERE slug = ?
		  AND status = 'published'
		LIMIT 1
	`, slug).Scan(
		&post.ID,
		&post.Title,
		&post.Slug,
		&post.Excerpt,
		&post.Content,
		&post.Status,
		&post.PublishedAt,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.BlogPost{}, sql.ErrNoRows
		}

		return models.BlogPost{}, fmt.Errorf(
			"get blog post %q: %w",
			slug,
			err,
		)
	}

	return post, nil
}

func (r *BlogRepository) ListAll(
	ctx context.Context,
) ([]models.BlogPost, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			id,
			title,
			slug,
			excerpt,
			content,
			status,
			published_at,
			created_at,
			updated_at
		FROM blog_posts
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("list blog posts: %w", err)
	}
	defer rows.Close()

	var posts []models.BlogPost

	for rows.Next() {
		var post models.BlogPost

		if err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Slug,
			&post.Excerpt,
			&post.Content,
			&post.Status,
			&post.PublishedAt,
			&post.CreatedAt,
			&post.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan blog post: %w", err)
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate blog posts: %w", err)
	}

	return posts, nil
}

func (r *BlogRepository) GetByID(
	ctx context.Context,
	id int64,
) (models.BlogPost, error) {
	var post models.BlogPost

	err := r.db.QueryRowContext(ctx, `
		SELECT
			id,
			title,
			slug,
			excerpt,
			content,
			status,
			published_at,
			created_at,
			updated_at
		FROM blog_posts
		WHERE id = ?
		LIMIT 1
	`, id).Scan(
		&post.ID,
		&post.Title,
		&post.Slug,
		&post.Excerpt,
		&post.Content,
		&post.Status,
		&post.PublishedAt,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.BlogPost{}, sql.ErrNoRows
		}

		return models.BlogPost{}, fmt.Errorf(
			"get blog post %d: %w",
			id,
			err,
		)
	}

	return post, nil
}

func (r *BlogRepository) Create(
	ctx context.Context,
	post models.BlogPost,
) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO blog_posts (
			title,
			slug,
			excerpt,
			content,
			status,
			published_at
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`,
		post.Title,
		post.Slug,
		post.Excerpt,
		post.Content,
		post.Status,
		post.PublishedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("create blog post: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf(
			"get created blog post id: %w",
			err,
		)
	}

	return id, nil
}

func (r *BlogRepository) Update(
	ctx context.Context,
	post models.BlogPost,
) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE blog_posts
		SET
			title = ?,
			slug = ?,
			excerpt = ?,
			content = ?,
			status = ?,
			published_at = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`,
		post.Title,
		post.Slug,
		post.Excerpt,
		post.Content,
		post.Status,
		post.PublishedAt,
		post.ID,
	)
	if err != nil {
		return fmt.Errorf(
			"update blog post %d: %w",
			post.ID,
			err,
		)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"get affected rows for blog post %d: %w",
			post.ID,
			err,
		)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *BlogRepository) Delete(
	ctx context.Context,
	id int64,
) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM blog_posts
		WHERE id = ?
	`, id)
	if err != nil {
		return fmt.Errorf(
			"delete blog post %d: %w",
			id,
			err,
		)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"get affected rows for blog post %d: %w",
			id,
			err,
		)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *BlogRepository) Publish(
	ctx context.Context,
	id int64,
) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE blog_posts
		SET
			status = 'published',
			published_at = COALESCE(
				published_at,
				CURRENT_TIMESTAMP
			),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, id)
	if err != nil {
		return fmt.Errorf(
			"publish blog post %d: %w",
			id,
			err,
		)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"get affected rows for blog post %d: %w",
			id,
			err,
		)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *BlogRepository) Unpublish(
	ctx context.Context,
	id int64,
) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE blog_posts
		SET
			status = 'draft',
			published_at = NULL,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, id)
	if err != nil {
		return fmt.Errorf(
			"unpublish blog post %d: %w",
			id,
			err,
		)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"get affected rows for blog post %d: %w",
			id,
			err,
		)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
