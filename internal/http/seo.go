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
	Location        string `xml:"loc"`
	LastModified    string `xml:"lastmod,omitempty"`
	ChangeFrequency string `xml:"changefreq,omitempty"`
	Priority        string `xml:"priority,omitempty"`
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
		{
			Location:        publicSiteURL + "/",
			ChangeFrequency: "weekly",
			Priority:        "1.0",
		},
		{
			Location:        publicSiteURL + "/work/portfolio",
			ChangeFrequency: "monthly",
			Priority:        "0.9",
		},
		{
			Location:        publicSiteURL + "/blog/",
			ChangeFrequency: "weekly",
			Priority:        "0.8",
		},
		{
			Location:        publicSiteURL + "/web-design-leeds/",
			ChangeFrequency: "monthly",
			Priority:        "0.9",
		},
	}

	for _, post := range posts {
		lastModified := post.UpdatedAt
		if lastModified.IsZero() && post.PublishedAt != nil {
			lastModified = *post.PublishedAt
		}

		entry := sitemapURL{
			Location:        publicSiteURL + "/blog/" + post.Slug,
			ChangeFrequency: "monthly",
			Priority:        "0.7",
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
