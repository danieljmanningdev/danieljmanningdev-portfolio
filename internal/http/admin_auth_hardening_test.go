package http

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	authservice "github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/auth"
)

func TestLogoutDoesNotClaimSuccessWhenRevocationFails(t *testing.T) {
	handler, _, _, db := newAdminAuthHTTPTest(t, false)
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	token := "test-session-token"
	form := url.Values{"csrf_token": {authservice.LogoutCSRFToken(token)}}
	req := httptest.NewRequest(http.MethodPost, logoutPath, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: adminSessionCookieName, Value: token})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusInternalServerError || rec.Header().Get("Location") != "" {
		t.Fatalf("logout response = %d", rec.Code)
	}
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == adminSessionCookieName && cookie.MaxAge < 0 {
			t.Fatal("cookie cleared despite failed revocation")
		}
	}
}

func TestLoginRejectsOversizedCredentialsBeforeDatabaseAccess(t *testing.T) {
	handler, _, _, db := newAdminAuthHTTPTest(t, false)
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{{strings.Repeat("a", 321), "password"}, {"admin@example.com", strings.Repeat("x", 73)}} {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, loginFormRequest(pair[0], pair[1]))
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", rec.Code)
		}
	}
}
