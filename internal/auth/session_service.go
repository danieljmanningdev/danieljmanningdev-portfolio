package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/models"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/repository"
)

const (
	adminSessionLifetime      = 24 * time.Hour
	adminSessionIdleTimeout   = 30 * time.Minute
	adminSessionTouchInterval = 5 * time.Minute
)

var (
	ErrInvalidSession = errors.New("invalid session")
	ErrExpiredSession = errors.New("expired session")
	ErrInactiveAdmin  = errors.New("admin account is inactive")
)

type SessionService struct {
	adminRepository   *repository.AdminRepository
	sessionRepository *repository.AdminSessionRepository
	now               func() time.Time
}

func NewSessionService(
	adminRepository *repository.AdminRepository,
	sessionRepository *repository.AdminSessionRepository,
) *SessionService {
	return &SessionService{
		adminRepository:   adminRepository,
		sessionRepository: sessionRepository,
		now: func() time.Time {
			return time.Now().UTC()
		},
	}
}

func (s *SessionService) CreateSession(
	ctx context.Context,
	adminID int64,
) (string, time.Time, error) {
	admin, err := s.adminRepository.GetByID(
		ctx,
		adminID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", time.Time{}, ErrInvalidSession
		}

		return "", time.Time{}, err
	}

	if !admin.IsActive {
		return "", time.Time{}, ErrInactiveAdmin
	}

	rawToken, err := GenerateSessionToken()
	if err != nil {
		return "", time.Time{}, err
	}

	tokenHash := HashSessionToken(rawToken)

	expiresAt := s.now().
		Add(adminSessionLifetime).
		UTC()

	_, err = s.sessionRepository.Create(
		ctx,
		admin.ID,
		tokenHash,
		expiresAt,
		s.now().UTC(),
	)
	if err != nil {
		return "", time.Time{}, err
	}

	return rawToken, expiresAt, nil
}

func (s *SessionService) Authenticate(
	ctx context.Context,
	rawToken string,
) (models.Admin, error) {
	if rawToken == "" {
		return models.Admin{}, ErrInvalidSession
	}

	tokenHash := HashSessionToken(rawToken)

	session, err := s.sessionRepository.GetByTokenHash(
		ctx,
		tokenHash,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Admin{}, ErrInvalidSession
		}

		return models.Admin{}, err
	}

	now := s.now().UTC()

	if !now.Before(session.ExpiresAt) {
		_ = s.sessionRepository.DeleteByTokenHash(
			ctx,
			tokenHash,
		)

		return models.Admin{}, ErrExpiredSession
	}

	if now.Sub(session.LastSeenAt) >= adminSessionIdleTimeout {
		_ = s.sessionRepository.DeleteByTokenHash(
			ctx,
			tokenHash,
		)

		return models.Admin{}, ErrExpiredSession
	}

	admin, err := s.adminRepository.GetByID(
		ctx,
		session.AdminID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = s.sessionRepository.DeleteByTokenHash(
				ctx,
				tokenHash,
			)

			return models.Admin{}, ErrInvalidSession
		}

		return models.Admin{}, err
	}

	if !admin.IsActive {
		_ = s.sessionRepository.DeleteByTokenHash(
			ctx,
			tokenHash,
		)

		return models.Admin{}, ErrInactiveAdmin
	}

	if now.Sub(session.LastSeenAt) >= adminSessionTouchInterval {
		if err := s.sessionRepository.Touch(
			ctx,
			session.ID,
			now,
		); err != nil {
			return models.Admin{}, err
		}
	}

	return admin, nil
}

func (s *SessionService) RevokeSession(
	ctx context.Context,
	rawToken string,
) error {
	if rawToken == "" {
		return nil
	}

	return s.sessionRepository.DeleteByTokenHash(
		ctx,
		HashSessionToken(rawToken),
	)
}

func (s *SessionService) RevokeAllAdminSessions(
	ctx context.Context,
	adminID int64,
) error {
	return s.sessionRepository.DeleteByAdminID(
		ctx,
		adminID,
	)
}

func (s *SessionService) DeleteExpiredSessions(
	ctx context.Context,
) error {
	return s.sessionRepository.DeleteExpired(
		ctx,
		s.now().UTC(),
	)
}
