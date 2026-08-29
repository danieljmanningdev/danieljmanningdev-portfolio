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
