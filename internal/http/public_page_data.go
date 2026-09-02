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

	personStructuredDataID              = publicSiteURL + "/#person"
	websiteStructuredDataID             = publicSiteURL + "/#website"
	professionalServiceStructuredDataID = publicSiteURL + "/#professional-service"
)

type breadcrumbLink struct {
	Name string
	URL  string
}

type publicPageData struct {
	Title          string
	Description    string
	OGTitle        string
	OGType         string
	CanonicalURL   string
	OGImage        string
	StructuredData template.JS
	CSPNonce       string
	RelatedLinks   []relatedLink
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
		canonicalURL = absolutePublicURL(path)
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

func (data publicPageData) withRelatedLinks(
	links ...relatedLink,
) publicPageData {
	data.RelatedLinks = append(
		[]relatedLink(nil),
		links...,
	)
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

func absolutePublicURL(path string) string {
	if path == "" || path == "/" {
		return publicSiteURL + "/"
	}

	return publicSiteURL + path
}

func personStructuredData() schema.Person {
	person := schema.NewPerson(
		"Daniel J. Manning",
		schema.WithID(personStructuredDataID),
	)

	person.URL = publicSiteURL + "/"
	person.JobTitle = "Digital Product Designer & Engineer"
	person.SocialProfiles = []schema.SocialProfile{
		"https://github.com/danieljmanningdev",
	}

	return person
}

func personStructuredDataReference() map[string]any {
	return map[string]any{
		"@id":   personStructuredDataID,
		"@type": "Person",
		"name":  "Daniel J. Manning",
		"url":   publicSiteURL + "/",
	}
}

func websiteStructuredData() schema.WebSite {
	website := schema.NewWebsite(
		"Daniel J. Manning",
		publicSiteURL+"/",
		schema.WithID(websiteStructuredDataID),
	)

	website.Publisher = personStructuredData().Node.Reference()

	return website
}

func professionalServiceStructuredData() schema.ProfessionalService {
	service := schema.NewProfessionalService(
		"Daniel J. Manning — Digital Product Design & Engineering",
		schema.WithID(professionalServiceStructuredDataID),
	)

	service.URL = publicSiteURL + "/"
	service.Description = "Digital product design, UI/UX, web design and Go software engineering for businesses in Leeds and across the United Kingdom."
	service.Image = defaultOGImage
	service.AreaServed = []string{
		"Leeds",
		"United Kingdom",
	}
	service.SameAs = []string{
		"https://github.com/danieljmanningdev",
	}

	return service
}

func homeStructuredData() schema.Graph {
	return schema.NewGraph(
		personStructuredData(),
		websiteStructuredData(),
		professionalServiceStructuredData(),
	)
}

func serviceStructuredData(
	name string,
	path string,
	description string,
	serviceType string,
	areaServed ...string,
) schema.Service {
	pageURL := absolutePublicURL(path)
	service := schema.NewService(
		name,
		schema.WithID(pageURL+"#service"),
	)

	service.URL = pageURL
	service.Description = description
	service.ServiceType = serviceType
	service.AreaServed = areaServed

	return service
}

func breadcrumbStructuredData(
	pageURL string,
	links ...breadcrumbLink,
) schema.BreadcrumbList {
	items := make(
		[]schema.ListItem,
		0,
		len(links),
	)

	for index, link := range links {
		items = append(
			items,
			schema.NewListItem(
				index+1,
				link.Name,
				link.URL,
			),
		)
	}

	return schema.NewBreadcrumbList(
		items,
		schema.WithID(pageURL+"#breadcrumb"),
	)
}

func servicePageStructuredData(
	name string,
	path string,
	description string,
	serviceType string,
	areaServed ...string,
) schema.Graph {
	pageURL := absolutePublicURL(path)

	return schema.NewGraph(
		serviceStructuredData(
			name,
			path,
			description,
			serviceType,
			areaServed...,
		),
		breadcrumbStructuredData(
			pageURL,
			breadcrumbLink{
				Name: "Home",
				URL:  publicSiteURL + "/",
			},
			breadcrumbLink{
				Name: name,
				URL:  pageURL,
			},
		),
	)
}

func caseStudyStructuredData(
	name string,
	path string,
	description string,
	keywords ...string,
) schema.Graph {
	pageURL := absolutePublicURL(path)
	creativeWork := map[string]any{
		"@id":         pageURL + "#case-study",
		"@type":       "CreativeWork",
		"name":        name,
		"description": description,
		"url":         pageURL,
		"image":       defaultOGImage,
		"author":      personStructuredDataReference(),
		"inLanguage":  "en-GB",
	}

	if len(keywords) > 0 {
		creativeWork["keywords"] = keywords
	}

	return schema.NewGraph(
		creativeWork,
		breadcrumbStructuredData(
			pageURL,
			breadcrumbLink{
				Name: "Home",
				URL:  publicSiteURL + "/",
			},
			breadcrumbLink{
				Name: name,
				URL:  pageURL,
			},
		),
	)
}

func blogStructuredData(
	title string,
	description string,
) schema.Graph {
	blogURL := publicSiteURL + "/blog/"
	blog := map[string]any{
		"@id":         blogURL + "#blog",
		"@type":       "Blog",
		"name":        title,
		"description": description,
		"url":         blogURL,
		"image":       defaultOGImage,
		"author":      personStructuredDataReference(),
		"inLanguage":  "en-GB",
	}

	return schema.NewGraph(
		blog,
		breadcrumbStructuredData(
			blogURL,
			breadcrumbLink{
				Name: "Home",
				URL:  publicSiteURL + "/",
			},
			breadcrumbLink{
				Name: "Journal",
				URL:  blogURL,
			},
		),
	)
}
