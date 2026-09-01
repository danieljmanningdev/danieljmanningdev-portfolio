// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package http

import (
	"context"
	"database/sql"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	authservice "github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/auth"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/database"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/repository"
)

func newAdminAuthHTTPTest(
	t *testing.T,
	secureCookies bool,
) (
	*AdminAuthHandler,
	*repository.AdminRepository,
	*repository.AdminSessionRepository,
	*sql.DB,
) {
	t.Helper()

	db, err := database.Open(
		context.Background(),
		":memory:",
	)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve test path")
	}

	migrationsDir := filepath.Join(
		filepath.Dir(filename),
		"..",
		"..",
		"migrations",
	)

	if err := database.RunMigrations(
		db.SQL,
		migrationsDir,
	); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	adminRepository :=
		repository.NewAdminRepository(db.SQL)

	sessionRepository :=
		repository.NewAdminSessionRepository(db.SQL)

	sessionService :=
		authservice.NewSessionService(
			adminRepository,
			sessionRepository,
		)

	loginTemplate := template.Must(
		template.New("base").Parse(`
			{{define "base"}}
				{{.Title}}
				{{.Email}}
				{{.Error}}
			{{end}}
		`),
	)

	dummyHash, err := authservice.HashPassword(
		"invalid-login-placeholder",
	)
	if err != nil {
		t.Fatalf("create dummy password hash: %v", err)
	}

	handler := &AdminAuthHandler{
		adminRepository: adminRepository,
		sessionService:  sessionService,
		loginLimiter:    authservice.NewLoginLimiter(),
		loginTemplates:  loginTemplate,
		secureCookies:   secureCookies,
		dummyHash:       dummyHash,
	}

	return handler,
		adminRepository,
		sessionRepository,
		db.SQL
}

func createAdminAuthTestAdmin(
	t *testing.T,
	adminRepository *repository.AdminRepository,
	email string,
	password string,
) int64 {
	t.Helper()

	passwordHash, err :=
		authservice.HashPassword(password)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	id, err := adminRepository.Create(
		context.Background(),
		email,
		passwordHash,
		"Daniel Manning",
	)
	if err != nil {
		t.Fatalf("create admin: %v", err)
	}

	return id
}

func loginFormRequest(
	email string,
	password string,
) *http.Request {
	csrfToken, err := authservice.GenerateCSRFToken()
	if err != nil {
		panic("generate test CSRF token: " + err.Error())
	}

	values := url.Values{}

	values.Set("email", email)
	values.Set("password", password)
	values.Set("csrf_token", csrfToken)

	req := httptest.NewRequest(
		http.MethodPost,
		loginPath,
		strings.NewReader(
			values.Encode(),
		),
	)

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	req.AddCookie(
		&http.Cookie{
			Name:  adminLoginCSRFCookie,
			Value: csrfToken,
		},
	)

	return req
}

func TestAdminLoginPage(t *testing.T) {
	handler, _, _, _ :=
		newAdminAuthHTTPTest(
			t,
			false,
		)

	req := httptest.NewRequest(
		http.MethodGet,
		loginPath,
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected 200, got %d",
			rec.Code,
		)
	}

	if !strings.Contains(
		rec.Body.String(),
		"Admin Login",
	) {
		t.Fatalf(
			"unexpected login page %q",
			rec.Body.String(),
		)
	}

	if cacheControl := rec.Header().Get(
		"Cache-Control",
	); cacheControl != "no-store" {
		t.Fatalf(
			"expected Cache-Control no-store, got %q",
			cacheControl,
		)
	}
}

func TestAdminLoginRejectsMissingCSRF(
	t *testing.T,
) {
	handler,
		adminRepository,
		_,
		_ := newAdminAuthHTTPTest(
		t,
		false,
	)

	createAdminAuthTestAdmin(
		t,
		adminRepository,
		"admin@example.com",
		"correct-password",
	)

	values := url.Values{}

	values.Set(
		"email",
		"admin@example.com",
	)

	values.Set(
		"password",
		"correct-password",
	)

	req := httptest.NewRequest(
		http.MethodPost,
		loginPath,
		strings.NewReader(
			values.Encode(),
		),
	)

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf(
			"expected 403, got %d",
			rec.Code,
		)
	}
}

func TestAdminLoginRejectsIncorrectCSRF(
	t *testing.T,
) {
	handler,
		adminRepository,
		_,
		_ := newAdminAuthHTTPTest(
		t,
		false,
	)

	createAdminAuthTestAdmin(
		t,
		adminRepository,
		"admin@example.com",
		"correct-password",
	)

	values := url.Values{}

	values.Set(
		"email",
		"admin@example.com",
	)

	values.Set(
		"password",
		"correct-password",
	)

	values.Set(
		"csrf_token",
		"submitted-token",
	)

	req := httptest.NewRequest(
		http.MethodPost,
		loginPath,
		strings.NewReader(
			values.Encode(),
		),
	)

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	req.AddCookie(
		&http.Cookie{
			Name:  adminLoginCSRFCookie,
			Value: "different-cookie-token",
		},
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf(
			"expected 403, got %d",
			rec.Code,
		)
	}
}

func TestAdminLoginAcceptsMatchingCSRF(
	t *testing.T,
) {
	handler,
		adminRepository,
		_,
		_ := newAdminAuthHTTPTest(
		t,
		false,
	)

	createAdminAuthTestAdmin(
		t,
		adminRepository,
		"admin@example.com",
		"correct-password",
	)

	req := loginFormRequest(
		"admin@example.com",
		"correct-password",
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf(
			"expected 303, got %d: %s",
			rec.Code,
			rec.Body.String(),
		)
	}

	if location := rec.Header().Get(
		"Location",
	); location != "/dashboard/" {
		t.Fatalf(
			"expected dashboard redirect, got %q",
			location,
		)
	}
}

func TestAdminLoginSuccess(t *testing.T) {
	handler,
		adminRepository,
		sessionRepository,
		_ := newAdminAuthHTTPTest(
		t,
		false,
	)

	adminID := createAdminAuthTestAdmin(
		t,
		adminRepository,
		"admin@example.com",
		"correct-password",
	)

	req := loginFormRequest(
		" ADMIN@example.com ",
		"correct-password",
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf(
			"expected 303, got %d: %s",
			rec.Code,
			rec.Body.String(),
		)
	}

	if location := rec.Header().Get(
		"Location",
	); location != "/dashboard/" {
		t.Fatalf(
			"expected dashboard redirect, got %q",
			location,
		)
	}

	result := rec.Result()

	var sessionCookie *http.Cookie

	for _, cookie := range result.Cookies() {
		if cookie.Name ==
			adminSessionCookieName {
			sessionCookie = cookie
			break
		}
	}

	if sessionCookie == nil {
		t.Fatal("expected admin session cookie")
	}

	if sessionCookie.Value == "" {
		t.Fatal(
			"expected non-empty session cookie",
		)
	}

	if !sessionCookie.HttpOnly {
		t.Fatal(
			"expected session cookie HttpOnly",
		)
	}

	if sessionCookie.Secure {
		t.Fatal(
			"expected local test cookie not to be Secure",
		)
	}

	if sessionCookie.SameSite !=
		http.SameSiteLaxMode {
		t.Fatalf(
			"expected SameSite=Lax, got %v",
			sessionCookie.SameSite,
		)
	}

	storedSession, err :=
		sessionRepository.GetByTokenHash(
			context.Background(),
			authservice.HashSessionToken(
				sessionCookie.Value,
			),
		)
	if err != nil {
		t.Fatalf(
			"get stored session: %v",
			err,
		)
	}

	if storedSession.AdminID != adminID {
		t.Fatalf(
			"expected admin ID %d, got %d",
			adminID,
			storedSession.AdminID,
		)
	}
}

func TestAdminLoginProductionCookieSecure(
	t *testing.T,
) {
	handler,
		adminRepository,
		_,
		_ := newAdminAuthHTTPTest(
		t,
		true,
	)

	createAdminAuthTestAdmin(
		t,
		adminRepository,
		"admin@example.com",
		"correct-password",
	)

	req := loginFormRequest(
		"admin@example.com",
		"correct-password",
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf(
			"expected 303, got %d",
			rec.Code,
		)
	}

	cookies := rec.Result().Cookies()

	if len(cookies) == 0 {
		t.Fatal("expected session cookie")
	}

	if !cookies[0].Secure {
		t.Fatal(
			"expected production cookie to be Secure",
		)
	}
}

func TestAdminLoginWrongPasswordUsesGenericError(
	t *testing.T,
) {
	handler,
		adminRepository,
		_,
		_ := newAdminAuthHTTPTest(
		t,
		false,
	)

	createAdminAuthTestAdmin(
		t,
		adminRepository,
		"admin@example.com",
		"correct-password",
	)

	req := loginFormRequest(
		"admin@example.com",
		"wrong-password",
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected 401, got %d",
			rec.Code,
		)
	}

	if !strings.Contains(
		rec.Body.String(),
		"Invalid email or password.",
	) {
		t.Fatalf(
			"unexpected body %q",
			rec.Body.String(),
		)
	}
}

func TestAdminLoginUnknownEmailUsesGenericError(
	t *testing.T,
) {
	handler, _, _, _ :=
		newAdminAuthHTTPTest(
			t,
			false,
		)

	req := loginFormRequest(
		"missing@example.com",
		"some-password",
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected 401, got %d",
			rec.Code,
		)
	}

	if !strings.Contains(
		rec.Body.String(),
		"Invalid email or password.",
	) {
		t.Fatalf(
			"unexpected body %q",
			rec.Body.String(),
		)
	}
}

func TestAdminLoginInactiveAdminUsesGenericError(
	t *testing.T,
) {
	handler,
		adminRepository,
		_,
		_ := newAdminAuthHTTPTest(
		t,
		false,
	)

	adminID := createAdminAuthTestAdmin(
		t,
		adminRepository,
		"admin@example.com",
		"correct-password",
	)

	if err := adminRepository.SetActive(
		context.Background(),
		adminID,
		false,
	); err != nil {
		t.Fatalf(
			"deactivate admin: %v",
			err,
		)
	}

	req := loginFormRequest(
		"admin@example.com",
		"correct-password",
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected 401, got %d",
			rec.Code,
		)
	}

	if !strings.Contains(
		rec.Body.String(),
		"Invalid email or password.",
	) {
		t.Fatalf(
			"unexpected body %q",
			rec.Body.String(),
		)
	}
}

func TestAdminLoginThrottlesRepeatedFailures(
	t *testing.T,
) {
	handler,
		adminRepository,
		_,
		_ := newAdminAuthHTTPTest(
		t,
		false,
	)

	createAdminAuthTestAdmin(
		t,
		adminRepository,
		"admin@example.com",
		"correct-password",
	)

	for i := 0; i < 5; i++ {
		req := loginFormRequest(
			"admin@example.com",
			"wrong-password",
		)

		req.RemoteAddr =
			"192.0.2.10:12345"

		rec := httptest.NewRecorder()

		handler.ServeHTTP(
			rec,
			req,
		)

		if rec.Code !=
			http.StatusUnauthorized {
			t.Fatalf(
				"attempt %d: expected 401, got %d",
				i+1,
				rec.Code,
			)
		}
	}

	req := loginFormRequest(
		"admin@example.com",
		"correct-password",
	)

	req.RemoteAddr =
		"192.0.2.10:12345"

	rec := httptest.NewRecorder()

	handler.ServeHTTP(
		rec,
		req,
	)

	if rec.Code !=
		http.StatusTooManyRequests {
		t.Fatalf(
			"expected 429 after repeated failures, got %d",
			rec.Code,
		)
	}

	if rec.Header().Get(
		"Retry-After",
	) == "" {
		t.Fatal(
			"expected Retry-After header",
		)
	}

	if !strings.Contains(
		rec.Body.String(),
		"Too many login attempts.",
	) {
		t.Fatalf(
			"unexpected throttled response %q",
			rec.Body.String(),
		)
	}
}

func TestAdminLoginThrottleIsScopedByIP(
	t *testing.T,
) {
	handler,
		adminRepository,
		_,
		_ := newAdminAuthHTTPTest(
		t,
		false,
	)

	createAdminAuthTestAdmin(
		t,
		adminRepository,
		"admin@example.com",
		"correct-password",
	)

	for i := 0; i < 5; i++ {
		req := loginFormRequest(
			"admin@example.com",
			"wrong-password",
		)

		req.RemoteAddr =
			"192.0.2.10:12345"

		rec := httptest.NewRecorder()

		handler.ServeHTTP(
			rec,
			req,
		)

		if rec.Code !=
			http.StatusUnauthorized {
			t.Fatalf(
				"attempt %d: expected 401, got %d",
				i+1,
				rec.Code,
			)
		}
	}

	req := loginFormRequest(
		"admin@example.com",
		"correct-password",
	)

	req.RemoteAddr =
		"192.0.2.11:12345"

	rec := httptest.NewRecorder()

	handler.ServeHTTP(
		rec,
		req,
	)

	if rec.Code !=
		http.StatusSeeOther {
		t.Fatalf(
			"expected different IP to remain allowed, got %d: %s",
			rec.Code,
			rec.Body.String(),
		)
	}
}

func TestAdminLogoutRevokesSession(
	t *testing.T,
) {
	handler,
		adminRepository,
		sessionRepository,
		_ := newAdminAuthHTTPTest(
		t,
		false,
	)

	adminID := createAdminAuthTestAdmin(
		t,
		adminRepository,
		"admin@example.com",
		"correct-password",
	)

	rawToken, _, err :=
		handler.sessionService.CreateSession(
			context.Background(),
			adminID,
		)
	if err != nil {
		t.Fatalf(
			"create session: %v",
			err,
		)
	}

	values := url.Values{}

	values.Set(
		"csrf_token",
		authservice.LogoutCSRFToken(
			rawToken,
		),
	)

	req := httptest.NewRequest(
		http.MethodPost,
		logoutPath,
		strings.NewReader(
			values.Encode(),
		),
	)

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	req.AddCookie(
		&http.Cookie{
			Name:  adminSessionCookieName,
			Value: rawToken,
		},
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf(
			"expected 303, got %d",
			rec.Code,
		)
	}

	if location := rec.Header().Get(
		"Location",
	); location != loginPath {
		t.Fatalf(
			"expected login redirect, got %q",
			location,
		)
	}

	_, err = sessionRepository.GetByTokenHash(
		context.Background(),
		authservice.HashSessionToken(
			rawToken,
		),
	)

	if err == nil {
		t.Fatal(
			"expected session to be revoked",
		)
	}

	cookies := rec.Result().Cookies()

	foundClearedCookie := false

	for _, cookie := range cookies {
		if cookie.Name ==
			adminSessionCookieName &&
			cookie.MaxAge < 0 {
			foundClearedCookie = true
		}
	}

	if !foundClearedCookie {
		t.Fatal(
			"expected browser session cookie to be cleared",
		)
	}
}

func TestAdminLogoutRejectsMissingCSRF(
	t *testing.T,
) {
	handler,
		adminRepository,
		sessionRepository,
		_ := newAdminAuthHTTPTest(
		t,
		false,
	)

	adminID := createAdminAuthTestAdmin(
		t,
		adminRepository,
		"admin@example.com",
		"correct-password",
	)

	rawToken, _, err :=
		handler.sessionService.CreateSession(
			context.Background(),
			adminID,
		)
	if err != nil {
		t.Fatalf(
			"create session: %v",
			err,
		)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		logoutPath,
		nil,
	)

	req.AddCookie(
		&http.Cookie{
			Name:  adminSessionCookieName,
			Value: rawToken,
		},
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf(
			"expected 403, got %d",
			rec.Code,
		)
	}

	_, err = sessionRepository.GetByTokenHash(
		context.Background(),
		authservice.HashSessionToken(
			rawToken,
		),
	)
	if err != nil {
		t.Fatalf(
			"expected session to remain after rejected logout: %v",
			err,
		)
	}
}

func TestAdminLogoutRejectsIncorrectCSRF(
	t *testing.T,
) {
	handler,
		adminRepository,
		sessionRepository,
		_ := newAdminAuthHTTPTest(
		t,
		false,
	)

	adminID := createAdminAuthTestAdmin(
		t,
		adminRepository,
		"admin@example.com",
		"correct-password",
	)

	rawToken, _, err :=
		handler.sessionService.CreateSession(
			context.Background(),
			adminID,
		)
	if err != nil {
		t.Fatalf(
			"create session: %v",
			err,
		)
	}

	values := url.Values{}

	values.Set(
		"csrf_token",
		"incorrect-token",
	)

	req := httptest.NewRequest(
		http.MethodPost,
		logoutPath,
		strings.NewReader(
			values.Encode(),
		),
	)

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	req.AddCookie(
		&http.Cookie{
			Name:  adminSessionCookieName,
			Value: rawToken,
		},
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf(
			"expected 403, got %d",
			rec.Code,
		)
	}

	_, err = sessionRepository.GetByTokenHash(
		context.Background(),
		authservice.HashSessionToken(
			rawToken,
		),
	)
	if err != nil {
		t.Fatalf(
			"expected session to remain after rejected logout: %v",
			err,
		)
	}
}

func TestAdminLogoutRejectsGET(
	t *testing.T,
) {
	handler, _, _, _ :=
		newAdminAuthHTTPTest(
			t,
			false,
		)

	req := httptest.NewRequest(
		http.MethodGet,
		logoutPath,
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code !=
		http.StatusMethodNotAllowed {
		t.Fatalf(
			"expected 405, got %d",
			rec.Code,
		)
	}
}

func TestRequireAdminRedirectsWithoutCookie(
	t *testing.T,
) {
	handler, _, _, _ :=
		newAdminAuthHTTPTest(
			t,
			false,
		)

	protected := RequireAdmin(
		handler.sessionService,
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				w.WriteHeader(
					http.StatusOK,
				)
			},
		),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/dashboard/",
		nil,
	)

	rec := httptest.NewRecorder()

	protected.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf(
			"expected 303, got %d",
			rec.Code,
		)
	}

	if location := rec.Header().Get(
		"Location",
	); location != loginPath {
		t.Fatalf(
			"expected login redirect, got %q",
			location,
		)
	}
}

func TestRequireAdminProvidesAdminContext(
	t *testing.T,
) {
	handler,
		adminRepository,
		_,
		_ := newAdminAuthHTTPTest(
		t,
		false,
	)

	adminID := createAdminAuthTestAdmin(
		t,
		adminRepository,
		"admin@example.com",
		"correct-password",
	)

	rawToken, _, err :=
		handler.sessionService.CreateSession(
			context.Background(),
			adminID,
		)
	if err != nil {
		t.Fatalf(
			"create session: %v",
			err,
		)
	}

	protected := RequireAdmin(
		handler.sessionService,
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				admin, ok := AdminFromContext(
					r.Context(),
				)
				if !ok {
					t.Fatal(
						"expected admin in context",
					)
				}

				if admin.ID != adminID {
					t.Fatalf(
						"expected admin ID %d, got %d",
						adminID,
						admin.ID,
					)
				}

				logoutCSRFToken, ok :=
					AdminLogoutCSRFTokenFromContext(
						r.Context(),
					)
				if !ok {
					t.Fatal(
						"expected logout CSRF token in context",
					)
				}

				expectedCSRFToken :=
					authservice.LogoutCSRFToken(
						rawToken,
					)

				if logoutCSRFToken !=
					expectedCSRFToken {
					t.Fatal(
						"unexpected logout CSRF token",
					)
				}

				w.WriteHeader(
					http.StatusOK,
				)
			},
		),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/dashboard/",
		nil,
	)

	req.AddCookie(
		&http.Cookie{
			Name:  adminSessionCookieName,
			Value: rawToken,
		},
	)

	rec := httptest.NewRecorder()

	protected.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected 200, got %d",
			rec.Code,
		)
	}
}

func TestRequireAdminRejectsExpiredSession(
	t *testing.T,
) {
	handler,
		adminRepository,
		sessionRepository,
		_ := newAdminAuthHTTPTest(
		t,
		false,
	)

	adminID := createAdminAuthTestAdmin(
		t,
		adminRepository,
		"admin@example.com",
		"correct-password",
	)

	rawToken := "expired-test-token"

	_, err := sessionRepository.Create(
		context.Background(),
		adminID,
		authservice.HashSessionToken(
			rawToken,
		),
		time.Now().UTC().Add(-time.Hour),
		time.Now().UTC().Add(-2*time.Hour),
	)
	if err != nil {
		t.Fatalf(
			"create expired session: %v",
			err,
		)
	}

	protected := RequireAdmin(
		handler.sessionService,
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				t.Fatal(
					"expired session reached protected handler",
				)
			},
		),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/dashboard/",
		nil,
	)

	req.AddCookie(
		&http.Cookie{
			Name:  adminSessionCookieName,
			Value: rawToken,
		},
	)

	rec := httptest.NewRecorder()

	protected.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf(
			"expected 303, got %d",
			rec.Code,
		)
	}

	if location := rec.Header().Get(
		"Location",
	); location != loginPath {
		t.Fatalf(
			"expected login redirect, got %q",
			location,
		)
	}
}
