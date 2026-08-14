package rendering

import (
	"html/template"
	"path/filepath"
)

func LoadPageTemplate(
	templateDir string,
	page string,
) (*template.Template, error) {
	return template.New("base").ParseFiles(
		filepath.Join(
			templateDir,
			"layouts",
			"base.html",
		),
		filepath.Join(
			templateDir,
			"components",
			"header.html",
		),
		filepath.Join(
			templateDir,
			"components",
			"footer.html",
		),
		filepath.Join(
			templateDir,
			"pages",
			page,
		),
	)
}
