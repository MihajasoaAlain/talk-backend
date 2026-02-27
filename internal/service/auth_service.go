package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"talk-backend/internal/models"
	"talk-backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthConfig struct {
	JWTSecret      string
	AccessTTL      time.Duration
	RefreshTTL     time.Duration
	MaxFailedLogin int
	LockDuration   time.Duration
	Issuer         string
}

type AuthService struct {
	users  repository.UserRepository
	tokens repository.RefreshTokenRepository
	audit  repository.AuditRepository
	cfg    AuthConfig
}

func NewAuthService(
	users repository.UserRepository,
	tokens repository.RefreshTokenRepository,
	audit repository.AuditRepository,
	cfg AuthConfig,
) *AuthService {
	return &AuthService{users: users, tokens: tokens, audit: audit, cfg: cfg}
}

func (s *AuthService) Register(username, email, password, avatarURL string) (*models.User, error) {
	_, err := s.users.FindByEmail(email)
	if err == nil {
		return nil, repository.ErrEmailAlreadyExists
	}
	if !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}

	u := &models.User{Username: username, Email: email, Password: password, AvatarURL: avatarURL}
	if err := s.users.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *AuthService) Login(email, password, ip, ua string) (accessToken string, refreshToken string, err error) {
	// Trouver user
	u, findErr := s.users.FindByEmail(email)
	if findErr != nil {
		s.auditLogin(nil, email, ip, ua, "login_fail")
		return "", "", ErrInvalidCredentials
	}

	// lock check
	if u.LockedUntil != nil && u.LockedUntil.After(time.Now()) {
		s.auditLogin(&u.ID, email, ip, ua, "login_fail")
		return "", "", ErrInvalidCredentials
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		u.FailedLoginAttempts++
		if u.FailedLoginAttempts >= s.cfg.MaxFailedLogin {
			lockUntil := time.Now().Add(s.cfg.LockDuration)
			u.LockedUntil = &lockUntil
			s.auditLogin(&u.ID, email, ip, ua, "account_locked")
		}
		_ = s.users.Update(u)
		s.auditLogin(&u.ID, email, ip, ua, "login_fail")
		return "", "", ErrInvalidCredentials
	}

	u.FailedLoginAttempts = 0
	u.LockedUntil = nil
	now := time.Now()
	u.LastLoginAt = &now
	_ = s.users.Update(u)

	s.auditLogin(&u.ID, email, ip, ua, "login_success")

	accessToken, err = s.signAccessToken(u.ID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = s.issueRefreshToken(u.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) Refresh(oldRefreshToken, ip, ua string) (newAccess string, newRefresh string, err error) {
	oldHash := hashToken(oldRefreshToken)

	rt, err := s.tokens.FindValidByHash(oldHash)
	if err != nil {
		s.auditLogin(nil, "", ip, ua, "refresh_fail")
		return "", "", ErrInvalidCredentials
	}

	now := time.Now()
	_ = s.tokens.Revoke(rt, now)

	newAccess, err = s.signAccessToken(rt.UserID)
	if err != nil {
		return "", "", err
	}

	newRefresh, err = s.issueRefreshToken(rt.UserID)
	if err != nil {
		return "", "", err
	}

	s.auditLogin(&rt.UserID, rt.User.Email, ip, ua, "refresh")
	return newAccess, newRefresh, nil
}

func (s *AuthService) Logout(refreshToken, ip, ua string) error {
	h := hashToken(refreshToken)
	rt, err := s.tokens.FindValidByHash(h)
	if err != nil {
		return nil // neutre
	}
	now := time.Now()
	_ = s.tokens.Revoke(rt, now)
	s.auditLogin(&rt.UserID, rt.User.Email, ip, ua, "logout")
	return nil
}

func (s *AuthService) signAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"iss": s.cfg.Issuer,
		"exp": time.Now().Add(s.cfg.AccessTTL).Unix(),
		"iat": time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *AuthService) issueRefreshToken(userID string) (string, error) {
	raw, err := randomToken(32)
	if err != nil {
		return "", err
	}
	rt := &models.RefreshToken{
		UserID:    userID,
		TokenHash: hashToken(raw),
		ExpiresAt: time.Now().Add(s.cfg.RefreshTTL),
	}
	if err := s.tokens.Create(rt); err != nil {
		return "", err
	}
	return raw, nil
}

func (s *AuthService) auditLogin(userID *string, email, ip, ua, event string) {
	_ = s.audit.Create(&models.AuditLog{
		UserID: userID,
		Event:  event,
		Email:  email,
		IP:     ip,
		UA:     ua,
	})
}

func randomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
