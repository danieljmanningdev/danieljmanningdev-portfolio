// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package http

import (
	"html/template"
	"net/http"

	"github.com/danieljmanningdev/go-jsonld-schema/schema"
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

	encoded, err := schema.Marshal(value)
	if err != nil {
		return ""
	}

	return template.JS(encoded)
}

func personStructuredData() schema.Person {
	person := schema.NewPerson("Daniel J. Manning")

	person.URL = publicSiteURL
	person.PictureURL = defaultOGImage
	person.JobTitle = "Digital Product Designer & Engineer"

	person.SocialProfiles = []schema.SocialProfile{
		"https://github.com/danieljmanningdev",
	}

	return person
}

func websiteStructuredData() schema.WebSite {
	website := schema.NewWebsite(
		"Daniel J. Manning",
		publicSiteURL,
	)

	return website
}

func homeStructuredData() schema.Graph {
	return schema.NewGraph(
		personStructuredData(),
		websiteStructuredData(),
	)
}
