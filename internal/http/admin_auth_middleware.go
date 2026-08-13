package http

import (
	"context"
	"errors"
	"net/http"

	authservice "github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/auth"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/models"
)

type adminContextKey struct{}

type adminLogoutCSRFContextKey struct{}

func AdminFromContext(
	ctx context.Context,
) (models.Admin, bool) {
	admin, ok := ctx.Value(
		adminContextKey{},
	).(models.Admin)

	return admin, ok
}

func AdminLogoutCSRFTokenFromContext(
	ctx context.Context,
) (string, bool) {
	token, ok := ctx.Value(
		adminLogoutCSRFContextKey{},
	).(string)

	return token, ok && token != ""
}

func RequireAdmin(
	sessionService *authservice.SessionService,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			cookie, err := r.Cookie(
				adminSessionCookieName,
			)
			if err != nil {
				redirectToAdminLogin(
					w,
					r,
				)
				return
			}

			admin, err := sessionService.Authenticate(
				r.Context(),
				cookie.Value,
			)
			if err != nil {
				switch {
				case errors.Is(
					err,
					authservice.ErrInvalidSession,
				),
					errors.Is(
						err,
						authservice.ErrExpiredSession,
					),
					errors.Is(
						err,
						authservice.ErrInactiveAdmin,
					):

					clearAdminSessionCookieForRequest(
						w,
						r,
					)

					redirectToAdminLogin(
						w,
						r,
					)

				default:
					http.Error(
						w,
						http.StatusText(
							http.StatusInternalServerError,
						),
						http.StatusInternalServerError,
					)
				}

				return
			}

			ctx := context.WithValue(
				r.Context(),
				adminContextKey{},
				admin,
			)

			ctx = context.WithValue(
				ctx,
				adminLogoutCSRFContextKey{},
				authservice.LogoutCSRFToken(
					cookie.Value,
				),
			)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		},
	)
}

func redirectToAdminLogin(
	w http.ResponseWriter,
	r *http.Request,
) {
	http.Redirect(
		w,
		r,
		loginPath,
		http.StatusSeeOther,
	)
}

func clearAdminSessionCookieForRequest(
	w http.ResponseWriter,
	r *http.Request,
) {
	secure := r.TLS != nil

	http.SetCookie(
		w,
		&http.Cookie{
			Name:     adminSessionCookieName,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   secure,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   -1,
		},
	)
}
