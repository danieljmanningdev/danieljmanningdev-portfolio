// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package blog

import (
	"regexp"
	"strings"
)

type FormErrors struct {
	Title   string
	Slug    string
	Excerpt string
	Content string
	Status  string
}

func (e FormErrors) Any() bool {
	return e.Title != "" ||
		e.Slug != "" ||
		e.Excerpt != "" ||
		e.Content != "" ||
		e.Status != ""
}

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func ValidateForm(form Form) FormErrors {
	var errors FormErrors

	title := strings.TrimSpace(form.Title)
	slug := strings.TrimSpace(form.Slug)
	excerpt := strings.TrimSpace(form.Excerpt)
	content := strings.TrimSpace(form.Content)
	status := strings.TrimSpace(form.Status)

	if title == "" {
		errors.Title = "Title is required."
	}

	if slug == "" {
		errors.Slug = "Slug is required."
	} else if !slugPattern.MatchString(slug) {
		errors.Slug = "Slug must use lowercase letters, numbers and hyphens only."
	}

	if excerpt == "" {
		errors.Excerpt = "Excerpt is required."
	}

	if content == "" {
		errors.Content = "Content is required."
	}

	switch status {
	case "draft", "published":
	default:
		errors.Status = "Status must be draft or published."
	}

	return errors
}
