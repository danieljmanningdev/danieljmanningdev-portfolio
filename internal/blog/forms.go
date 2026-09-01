// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package blog

import (
	"strings"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/models"
)

type Form struct {
	Title   string
	Slug    string
	Excerpt string
	Content string
	Status  string
}

func FormFromPost(post models.BlogPost) Form {
	return Form{
		Title:   post.Title,
		Slug:    post.Slug,
		Excerpt: post.Excerpt,
		Content: post.Content,
		Status:  post.Status,
	}
}

func (f Form) ToPost() models.BlogPost {
	return models.BlogPost{
		Title:   strings.TrimSpace(f.Title),
		Slug:    strings.TrimSpace(f.Slug),
		Excerpt: strings.TrimSpace(f.Excerpt),
		Content: strings.TrimSpace(f.Content),
		Status:  strings.TrimSpace(f.Status),
	}
}

func (f Form) ToPostWithID(id int64) models.BlogPost {
	post := f.ToPost()
	post.ID = id

	return post
}
