// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package rendering

import (
	"html/template"
	"path/filepath"
	"strings"
)

func LoadPageTemplate(
	templateDir string,
	page string,
) (*template.Template, error) {
	switch {
	case page == "admin/login.html":
		return loadAuthPageTemplate(
			templateDir,
			page,
		)

	case page == "blog.html",
		page == "blog-post.html",
		strings.HasPrefix(page, "public/"):
		return loadPublicPageTemplate(
			templateDir,
			page,
		)

	default:
		return loadAdminPageTemplate(
			templateDir,
			page,
		)
	}
}

func loadPublicPageTemplate(
	templateDir string,
	page string,
) (*template.Template, error) {
	return template.New("base").ParseFiles(
		filepath.Join(templateDir, "layouts", "base.html"),
		filepath.Join(templateDir, "components", "header.html"),
		filepath.Join(templateDir, "components", "footer.html"),
		filepath.Join(templateDir, "pages", page),
	)
}

func loadAdminPageTemplate(
	templateDir string,
	page string,
) (*template.Template, error) {
	return template.New("base").ParseFiles(
		filepath.Join(templateDir, "layouts", "admin-base.html"),
		filepath.Join(templateDir, "components", "admin-header.html"),
		filepath.Join(templateDir, "components", "admin-footer.html"),
		filepath.Join(templateDir, "pages", page),
	)
}

func loadAuthPageTemplate(
	templateDir string,
	page string,
) (*template.Template, error) {
	return template.New("base").ParseFiles(
		filepath.Join(templateDir, "layouts", "auth-base.html"),
		filepath.Join(templateDir, "pages", page),
	)
}
