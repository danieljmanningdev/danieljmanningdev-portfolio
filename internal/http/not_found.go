// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package http

import (
	"bytes"
	"html/template"
	"net/http"
)

type notFoundPageData struct {
	publicPageData
	Path string
}

func newNotFoundPageData(path string) notFoundPageData {
	return notFoundPageData{
		publicPageData: newPublicPageData(
			"Page Not Found — Daniel J. Manning",
			"The requested page could not be found.",
			"",
			"website",
			nil,
		),
		Path: path,
	}
}

func renderNotFoundPage(
	w http.ResponseWriter,
	tmpl *template.Template,
	path string,
) {
	var body bytes.Buffer

	if err := tmpl.ExecuteTemplate(
		&body,
		"base",
		newNotFoundPageData(path),
	); err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)
	w.Header().Set(
		"X-Robots-Tag",
		"noindex, nofollow",
	)
	w.WriteHeader(http.StatusNotFound)

	_, _ = body.WriteTo(w)
}
