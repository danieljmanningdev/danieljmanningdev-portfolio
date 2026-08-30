// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package http

import (
	"bytes"
	"database/sql"
	"errors"
	"html/template"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	authservice "github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/auth"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/rendering"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/repository"
)

const (
	adminSessionCookieName = "djm_admin_session"
	adminLoginCSRFCookie   = "djm_login_csrf"
	loginPath              = "/login"
	logoutPath             = "/logout"

	loginCSRFLifetime = 10 * time.Minute
)

type adminLoginPageData struct {
	publicPageData

	Email     string
	Error     string
	CSRFToken string
}

type AdminAuthHandler struct {
	adminRepository *repository.AdminRepository
	sessionService  *authservice.SessionService
	loginLimiter    *authservice.LoginLimiter
	loginTemplates  *template.Template
	secureCookies   bool
	dummyHash       string
}

func NewAdminAuthHandler(
	db *sql.DB,
	templateDir string,
	secureCookies bool,
) (*AdminAuthHandler, error) {
	adminRepository := repository.NewAdminRepository(db)

	sessionRepository :=
		repository.NewAdminSessionRepository(db)

	sessionService := authservice.NewSessionService(
		adminRepository,
		sessionRepository,
	)

	loginTemplates, err := rendering.LoadPageTemplate(
		templateDir,
		"admin/login.html",
	)
	if err != nil {
		return nil, err
	}

	/*
		The dummy bcrypt hash lets failed logins for unknown email
		addresses still perform a bcrypt comparison.

		This helps avoid making "email exists" noticeably cheaper
		than "wrong password".
	*/
	dummyHash, err := authservice.HashPassword(
		"invalid-login-placeholder",
	)
	if err != nil {
		return nil, err
	}

	return &AdminAuthHandler{
		adminRepository: adminRepository,
		sessionService:  sessionService,
		loginLimiter:    authservice.NewLoginLimiter(),
		loginTemplates:  loginTemplates,
		secureCookies:   secureCookies,
		dummyHash:       dummyHash,
	}, nil
}

func (h *AdminAuthHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	switch r.URL.Path {
	case loginPath:
		h.handleLogin(w, r)

	case logoutPath:
		h.handleLogout(w, r)

	default:
		http.NotFound(w, r)
	}
}

func (h *AdminAuthHandler) handleLogin(
	w http.ResponseWriter,
	r *http.Request,
) {
	switch r.Method {
	case http.MethodGet:
		h.showLogin(w, r)

	case http.MethodPost:
		h.submitLogin(w, r)

	default:
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
	}
}

func (h *AdminAuthHandler) showLogin(
	w http.ResponseWriter,
	r *http.Request,
) {
	if h.requestHasValidAdminSession(r) {
		http.Redirect(
			w,
			r,
			"/dashboard/",
			http.StatusSeeOther,
		)
		return
	}

	csrfToken, err := authservice.GenerateCSRFToken()
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	h.setLoginCSRFCookie(
		w,
		csrfToken,
	)

	h.renderLogin(
		w,
		adminLoginPageData{
			publicPageData: adminLoginMetadata(),
			CSRFToken:      csrfToken,
		},
		http.StatusOK,
	)
}

func (h *AdminAuthHandler) submitLogin(
	w http.ResponseWriter,
	r *http.Request,
) {
	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest,
		)
		return
	}

	csrfCookie, err := r.Cookie(
		adminLoginCSRFCookie,
	)
	if err != nil ||
		!authservice.VerifyCSRFToken(
			csrfCookie.Value,
			r.FormValue("csrf_token"),
		) {
		http.Error(
			w,
			http.StatusText(http.StatusForbidden),
			http.StatusForbidden,
		)
		return
	}

	email := strings.ToLower(
		strings.TrimSpace(
			r.FormValue("email"),
		),
	)

	password := r.FormValue("password")

	clientIP := loginClientIP(r)

	if allowed, retryAfter :=
		h.loginLimiter.Check(
			clientIP,
			email,
		); !allowed {
		h.renderLoginThrottled(
			w,
			retryAfter,
		)
		return
	}

	if email == "" || password == "" {
		h.rejectInvalidLogin(
			w,
			clientIP,
			email,
		)
		return
	}

	admin, err := h.adminRepository.GetByEmail(
		r.Context(),
		email,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			/*
				Still perform bcrypt work when the email does
				not exist so the two failure cases behave more
				similarly.
			*/
			_ = authservice.VerifyPassword(
				h.dummyHash,
				password,
			)

			h.rejectInvalidLogin(
				w,
				clientIP,
				email,
			)
			return
		}

		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	if !admin.IsActive ||
		!authservice.VerifyPassword(
			admin.PasswordHash,
			password,
		) {
		h.rejectInvalidLogin(
			w,
			clientIP,
			email,
		)
		return
	}

	/*
		Remove any existing browser session before issuing a new
		one. This also gives successful authentication a freshly
		generated session token.
	*/
	if existingCookie, err := r.Cookie(
		adminSessionCookieName,
	); err == nil {
		_ = h.sessionService.RevokeSession(
			r.Context(),
			existingCookie.Value,
		)
	}

	rawToken, expiresAt, err :=
		h.sessionService.CreateSession(
			r.Context(),
			admin.ID,
		)
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	h.loginLimiter.ResetCredential(
		clientIP,
		email,
	)

	h.setSessionCookie(
		w,
		rawToken,
		expiresAt,
	)

	h.clearLoginCSRFCookie(w)

	http.Redirect(
		w,
		r,
		"/dashboard/",
		http.StatusSeeOther,
	)
}

func (h *AdminAuthHandler) handleLogout(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest,
		)
		return
	}

	cookie, err := r.Cookie(
		adminSessionCookieName,
	)
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusForbidden),
			http.StatusForbidden,
		)
		return
	}

	expectedCSRFToken :=
		authservice.LogoutCSRFToken(
			cookie.Value,
		)

	if !authservice.VerifyCSRFToken(
		expectedCSRFToken,
		r.FormValue("csrf_token"),
	) {
		http.Error(
			w,
			http.StatusText(http.StatusForbidden),
			http.StatusForbidden,
		)
		return
	}

	_ = h.sessionService.RevokeSession(
		r.Context(),
		cookie.Value,
	)

	h.clearSessionCookie(w)

	http.Redirect(
		w,
		r,
		loginPath,
		http.StatusSeeOther,
	)
}

func (h *AdminAuthHandler) requestHasValidAdminSession(
	r *http.Request,
) bool {
	cookie, err := r.Cookie(
		adminSessionCookieName,
	)
	if err != nil {
		return false
	}

	_, err = h.sessionService.Authenticate(
		r.Context(),
		cookie.Value,
	)

	return err == nil
}

func (h *AdminAuthHandler) rejectInvalidLogin(
	w http.ResponseWriter,
	clientIP string,
	email string,
) {
	h.loginLimiter.RecordFailure(
		clientIP,
		email,
	)

	h.renderInvalidLogin(
		w,
		email,
	)
}

func (h *AdminAuthHandler) renderLoginThrottled(
	w http.ResponseWriter,
	retryAfter time.Duration,
) {
	retryAfterSeconds := int(
		math.Ceil(
			retryAfter.Seconds(),
		),
	)

	if retryAfterSeconds < 1 {
		retryAfterSeconds = 1
	}

	w.Header().Set(
		"Retry-After",
		strconv.Itoa(
			retryAfterSeconds,
		),
	)

	csrfToken, err := authservice.GenerateCSRFToken()
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	h.setLoginCSRFCookie(
		w,
		csrfToken,
	)

	h.renderLogin(
		w,
		adminLoginPageData{
			publicPageData: adminLoginMetadata(),
			Error:          "Too many login attempts. Please try again later.",
			CSRFToken:      csrfToken,
		},
		http.StatusTooManyRequests,
	)
}

func loginClientIP(
	r *http.Request,
) string {
	host, _, err := net.SplitHostPort(
		r.RemoteAddr,
	)
	if err == nil && host != "" {
		return host
	}

	if r.RemoteAddr != "" {
		return r.RemoteAddr
	}

	return "unknown"
}

func (h *AdminAuthHandler) renderInvalidLogin(
	w http.ResponseWriter,
	email string,
) {
	csrfToken, err := authservice.GenerateCSRFToken()
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	h.setLoginCSRFCookie(
		w,
		csrfToken,
	)

	h.renderLogin(
		w,
		adminLoginPageData{
			publicPageData: adminLoginMetadata(),
			Email:          email,
			Error:          "Invalid email or password.",
			CSRFToken:      csrfToken,
		},
		http.StatusUnauthorized,
	)
}

func adminLoginMetadata() publicPageData {
	return publicPageData{
		Title:       "Admin Login — Daniel J. Manning",
		Description: "Secure administrative access to the Daniel J. Manning client workspace.",
		OGTitle:     "Admin Login — Daniel J. Manning",
		OGType:      "website",
	}
}

func (h *AdminAuthHandler) renderLogin(
	w http.ResponseWriter,
	data adminLoginPageData,
	status int,
) {
	var body bytes.Buffer

	if err := h.loginTemplates.ExecuteTemplate(
		&body,
		"base",
		data,
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
		"Cache-Control",
		"no-store",
	)

	w.WriteHeader(status)

	_, _ = body.WriteTo(w)
}

func (h *AdminAuthHandler) setSessionCookie(
	w http.ResponseWriter,
	rawToken string,
	expiresAt time.Time,
) {
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     adminSessionCookieName,
			Value:    rawToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   h.secureCookies,
			SameSite: http.SameSiteLaxMode,
			Expires:  expiresAt,
		},
	)
}

func (h *AdminAuthHandler) clearSessionCookie(
	w http.ResponseWriter,
) {
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     adminSessionCookieName,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   h.secureCookies,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   -1,
			Expires:  time.Unix(1, 0).UTC(),
		},
	)
}

func (h *AdminAuthHandler) setLoginCSRFCookie(
	w http.ResponseWriter,
	token string,
) {
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     adminLoginCSRFCookie,
			Value:    token,
			Path:     loginPath,
			HttpOnly: true,
			Secure:   h.secureCookies,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   int(loginCSRFLifetime.Seconds()),
		},
	)
}

func (h *AdminAuthHandler) clearLoginCSRFCookie(
	w http.ResponseWriter,
) {
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     adminLoginCSRFCookie,
			Value:    "",
			Path:     loginPath,
			HttpOnly: true,
			Secure:   h.secureCookies,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   -1,
			Expires:  time.Unix(1, 0).UTC(),
		},
	)
}
