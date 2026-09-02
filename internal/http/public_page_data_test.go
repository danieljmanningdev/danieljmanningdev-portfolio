// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package http

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/models"
)

func TestHomeStructuredDataIncludesCoreEntities(
	t *testing.T,
) {
	document := decodeStructuredData(
		t,
		homeStructuredData(),
	)

	assertStructuredDataTypes(
		t,
		document,
		"Person",
		"WebSite",
		"ProfessionalService",
	)
}

func TestServicePageStructuredDataIncludesServiceAndBreadcrumbs(
	t *testing.T,
) {
	document := decodeStructuredData(
		t,
		servicePageStructuredData(
			"Web Design",
			"/web-design/",
			"Responsive web design.",
			"Web design",
			"United Kingdom",
		),
	)

	assertStructuredDataTypes(
		t,
		document,
		"Service",
		"BreadcrumbList",
	)

	service := structuredDataNodeByType(
		t,
		document,
		"Service",
	)

	if got := service["url"]; got != publicSiteURL+"/web-design/" {
		t.Fatalf(
			"expected canonical service URL, got %#v",
			got,
		)
	}
}

func TestCaseStudyStructuredDataIncludesCreativeWorkAndBreadcrumbs(
	t *testing.T,
) {
	document := decodeStructuredData(
		t,
		caseStudyStructuredData(
			"Salon Rebuild",
			"/work/salon-rebuild/",
			"A modern salon website concept.",
			"Go",
			"UI/UX design",
		),
	)

	assertStructuredDataTypes(
		t,
		document,
		"CreativeWork",
		"BreadcrumbList",
	)
}

func TestBlogPostStructuredDataIncludesDatesAndBreadcrumbs(
	t *testing.T,
) {
	publishedAt := time.Date(
		2026,
		time.September,
		1,
		9,
		30,
		0,
		0,
		time.UTC,
	)
	updatedAt := publishedAt.Add(2 * time.Hour)

	document := decodeStructuredData(
		t,
		blogPostStructuredData(
			models.BlogPost{
				Title:       "Example article",
				Slug:        "example-article",
				PublishedAt: &publishedAt,
				UpdatedAt:   updatedAt,
			},
			"Example description.",
			"/blog/example-article",
		),
	)

	assertStructuredDataTypes(
		t,
		document,
		"BlogPosting",
		"BreadcrumbList",
	)

	posting := structuredDataNodeByType(
		t,
		document,
		"BlogPosting",
	)

	if got := posting["datePublished"]; got != publishedAt.Format(time.RFC3339) {
		t.Fatalf(
			"expected publication date %q, got %#v",
			publishedAt.Format(time.RFC3339),
			got,
		)
	}

	if got := posting["dateModified"]; got != updatedAt.Format(time.RFC3339) {
		t.Fatalf(
			"expected modification date %q, got %#v",
			updatedAt.Format(time.RFC3339),
			got,
		)
	}
}

func TestRelatedLinksForSalonRetrospective(
	t *testing.T,
) {
	links := relatedLinksForJournalPost(
		"i-built-a-website-for-50p-an-hour-never-again",
	)

	if len(links) != 3 {
		t.Fatalf(
			"expected three related links, got %d",
			len(links),
		)
	}

	if links[0].URL != "/work/salon-rebuild/" {
		t.Fatalf(
			"expected Salon Rebuild to be the primary related link, got %q",
			links[0].URL,
		)
	}
}

func decodeStructuredData(
	t *testing.T,
	value any,
) map[string]any {
	t.Helper()

	encoded := marshalStructuredData(value)
	if encoded == "" {
		t.Fatal("expected structured data to marshal")
	}

	var document map[string]any
	if err := json.Unmarshal(
		[]byte(encoded),
		&document,
	); err != nil {
		t.Fatalf(
			"decode structured data: %v",
			err,
		)
	}

	if got := document["@context"]; got != "https://schema.org" {
		t.Fatalf(
			"expected schema.org context, got %#v",
			got,
		)
	}

	return document
}

func assertStructuredDataTypes(
	t *testing.T,
	document map[string]any,
	expectedTypes ...string,
) {
	t.Helper()

	for _, expectedType := range expectedTypes {
		_ = structuredDataNodeByType(
			t,
			document,
			expectedType,
		)
	}
}

func structuredDataNodeByType(
	t *testing.T,
	document map[string]any,
	expectedType string,
) map[string]any {
	t.Helper()

	graph, ok := document["@graph"].([]any)
	if !ok {
		t.Fatalf(
			"expected @graph array, got %#v",
			document["@graph"],
		)
	}

	for _, rawNode := range graph {
		node, ok := rawNode.(map[string]any)
		if !ok {
			continue
		}

		if node["@type"] == expectedType {
			return node
		}
	}

	t.Fatalf(
		"expected structured data type %q in %#v",
		expectedType,
		graph,
	)
	return nil
}
