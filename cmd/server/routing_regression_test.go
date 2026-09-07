package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	apphttp "github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/http"
)

func TestPublicRoutesAreExactAndOldPortfolioLinkRedirects(t *testing.T) {
	fallback := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.NotFound(w, r) })
	page := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	router := newRouter(routerDependencies{
		homeHandler: fallback, portfolioCaseStudyHandler: page,
		publicPageRoutes: []apphttp.PublicPageRoute{{Path: "/web-design/", Handler: page}},
		adminAuthHandler: page, dashboardHandler: page, clientsHandler: page,
		projectsHandler: page, contractsHandler: page, blogAdminHandler: page,
		sessionService: newRouterTestSessionService(t), blogHandler: &apphttp.BlogHandler{},
	})
	for _, tc := range []struct {
		path   string
		status int
	}{
		{"/web-design/", 200}, {"/web-design/does-not-exist", 404},
		{"/work/portfolio", 200}, {"/work/portfolio/", 308}, {"/work/portfolio/does-not-exist", 404},
	} {
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))
		if rec.Code != tc.status {
			t.Errorf("%s: got %d, want %d", tc.path, rec.Code, tc.status)
		}
		if tc.status == 308 && rec.Header().Get("Location") != "/work/portfolio" {
			t.Error("wrong canonical redirect")
		}
	}
}
