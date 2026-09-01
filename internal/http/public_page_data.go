// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package http

import (
	"encoding/json"
	"html/template"
	"net/http"
)

const (
	publicSiteURL  = "https://danieljmanningdev.com"
	defaultOGImage = publicSiteURL + "/static/images/og-card.png"
)

type publicPageData struct {
	Title          string
	Description    string
	OGTitle        string
	OGType         string
	CanonicalURL   string
	OGImage        string
	StructuredData template.JS
	CSPNonce       string
}

func newPublicPageData(
	title string,
	description string,
	path string,
	ogType string,
	structuredData any,
) publicPageData {
	canonicalURL := ""
	if path != "" {
		canonicalURL = publicSiteURL + path
	}

	return publicPageData{
		Title:          title,
		Description:    description,
		OGTitle:        title,
		OGType:         ogType,
		CanonicalURL:   canonicalURL,
		OGImage:        defaultOGImage,
		StructuredData: marshalStructuredData(structuredData),
	}
}

func (data publicPageData) withRequest(
	r *http.Request,
) publicPageData {
	data.CSPNonce = CSPNonceFromContext(r.Context())
	return data
}

func marshalStructuredData(value any) template.JS {
	if value == nil {
		return ""
	}

	encoded, err := json.Marshal(value)
	if err != nil {
		return ""
	}

	return template.JS(encoded)
}

func personStructuredData() map[string]any {
	return map[string]any{
		"@context": "https://schema.org",
		"@type":    "Person",
		"name":     "Daniel J. Manning",
		"url":      publicSiteURL,
		"image":    defaultOGImage,
		"jobTitle": "Digital Product Designer & Engineer",
		"sameAs": []string{
			"https://github.com/danieljmanningdev",
		},
		"knowsAbout": []string{
			"Digital product design",
			"User experience design",
			"Go software engineering",
			"HTMX",
			"Web application security",
		},
	}
}
