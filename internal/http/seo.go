// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package http

import (
	"encoding/xml"
	"net/http"
	"time"
)

type sitemapURLSet struct {
	XMLName xml.Name     `xml:"urlset"`
	XMLNS   string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

type sitemapURL struct {
	Location     string `xml:"loc"`
	LastModified string `xml:"lastmod,omitempty"`
}

func RobotsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set(
		"Content-Type",
		"text/plain; charset=utf-8",
	)
	w.Header().Set(
		"Cache-Control",
		"public, max-age=86400",
	)

	_, _ = w.Write([]byte(
		"User-agent: *\n" +
			"Allow: /\n" +
			"Disallow: /login\n" +
			"Disallow: /logout\n" +
			"Disallow: /dashboard/\n" +
			"Sitemap: " + publicSiteURL + "/sitemap.xml\n",
	))
}

func (h *BlogHandler) Sitemap(
	w http.ResponseWriter,
	r *http.Request,
) {
	posts, err := h.repository.ListPublished(r.Context())
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	urls := []sitemapURL{
		{Location: publicSiteURL + "/"},
		{Location: publicSiteURL + "/work/portfolio"},
		{Location: publicSiteURL + "/blog/"},
	}

	for _, page := range PublicPages {
		urls = append(
			urls,
			sitemapURL{
				Location: publicSiteURL + page.Path,
			},
		)
	}

	for _, post := range posts {
		lastModified := post.UpdatedAt
		if lastModified.IsZero() && post.PublishedAt != nil {
			lastModified = *post.PublishedAt
		}

		entry := sitemapURL{
			Location: publicSiteURL + "/blog/" + post.Slug,
		}

		if !lastModified.IsZero() {
			entry.LastModified = lastModified.UTC().Format(
				time.DateOnly,
			)
		}

		urls = append(urls, entry)
	}

	w.Header().Set(
		"Content-Type",
		"application/xml; charset=utf-8",
	)
	w.Header().Set(
		"Cache-Control",
		"public, max-age=3600",
	)
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte(xml.Header))

	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")

	if err := encoder.Encode(sitemapURLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}); err != nil {
		return
	}
}
