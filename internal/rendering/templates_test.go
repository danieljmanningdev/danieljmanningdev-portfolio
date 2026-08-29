package rendering

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestLoadPageTemplateParsesApplicationPages(
	t *testing.T,
) {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve test path")
	}

	templateDir := filepath.Join(
		filepath.Dir(filename),
		"..",
		"..",
		"web",
		"templates",
	)

	pages := []string{
		"public/home.html",
		"public/portfolio.html",
		"public/404.html",
		"blog.html",
		"blog-post.html",
		"admin/login.html",
		"admin/dashboard.html",
		"admin/activity.html",
		"admin/blog/index.html",
		"admin/blog/new.html",
		"admin/blog/edit.html",
		"admin/clients/index.html",
		"admin/clients/new.html",
		"admin/clients/show.html",
		"admin/clients/edit.html",
		"admin/projects/index.html",
		"admin/projects/new.html",
		"admin/projects/show.html",
		"admin/projects/edit.html",
		"admin/contracts/index.html",
		"admin/contracts/new.html",
		"admin/contracts/show.html",
		"admin/contracts/edit.html",
	}

	for _, page := range pages {
		t.Run(page, func(t *testing.T) {
			if _, err := LoadPageTemplate(
				templateDir,
				page,
			); err != nil {
				t.Fatalf(
					"parse template %q: %v",
					page,
					err,
				)
			}
		})
	}
}
