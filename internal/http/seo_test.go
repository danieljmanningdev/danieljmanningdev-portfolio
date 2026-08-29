package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/database"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/repository"
)

func TestRobotsHandler(
	t *testing.T,
) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodGet,
		"/robots.txt",
		nil,
	)

	RobotsHandler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	body := recorder.Body.String()
	for _, expected := range []string{
		"User-agent: *",
		"Allow: /",
		"Disallow: /dashboard/",
		"Sitemap: https://danieljmanningdev.com/sitemap.xml",
	} {
		if !strings.Contains(body, expected) {
			t.Fatalf(
				"expected robots.txt to contain %q, got %q",
				expected,
				body,
			)
		}
	}
}

func TestSitemapIncludesOnlyPublishedJournalPosts(
	t *testing.T,
) {
	ctx := context.Background()
	databaseHandle, err := database.Open(ctx, ":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer func() {
		_ = databaseHandle.Close()
	}()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve test path")
	}

	migrationsDirectory := filepath.Join(
		filepath.Dir(filename),
		"..",
		"..",
		"migrations",
	)

	if err := database.RunMigrations(
		databaseHandle.SQL,
		migrationsDirectory,
	); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	if _, err := databaseHandle.SQL.Exec(`
		INSERT INTO blog_posts (
			title,
			slug,
			excerpt,
			content,
			status,
			published_at,
			updated_at
		)
		VALUES
			(
				'Published Post',
				'published-post',
				'Published excerpt',
				'Published content',
				'published',
				'2026-08-20 12:00:00',
				'2026-08-21 12:00:00'
			),
			(
				'Draft Post',
				'draft-post',
				'Draft excerpt',
				'Draft content',
				'draft',
				NULL,
				'2026-08-22 12:00:00'
			)
	`); err != nil {
		t.Fatalf("insert posts: %v", err)
	}

	handler := &BlogHandler{
		repository: repository.NewBlogRepository(
			databaseHandle.SQL,
		),
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodGet,
		"/sitemap.xml",
		nil,
	)

	handler.Sitemap(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	body := recorder.Body.String()
	for _, expected := range []string{
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>",
		"https://danieljmanningdev.com/",
		"https://danieljmanningdev.com/work/portfolio",
		"https://danieljmanningdev.com/blog/",
		"https://danieljmanningdev.com/blog/published-post",
		"<lastmod>2026-08-21</lastmod>",
	} {
		if !strings.Contains(body, expected) {
			t.Fatalf(
				"expected sitemap to contain %q, got %q",
				expected,
				body,
			)
		}
	}

	if strings.Contains(body, "draft-post") {
		t.Fatalf("expected draft post to be excluded, got %q", body)
	}
}
