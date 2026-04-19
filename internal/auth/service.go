package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/andrew-hayworth22/critiquefi-service/internal/models"
	"github.com/andrew-hayworth22/critiquefi-service/internal/store"
	"github.com/andrew-hayworth22/critiquefi-service/pkg/crypto"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrDuplicate          = errors.New("duplicate record")
)

type Store interface {
	CreateUser(ctx context.Context, user models.NewUser) (id int64, err error)
	GetUserByID(ctx context.Context, id int64) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (user models.User, err error)
	CheckTakenUserFields(ctx context.Context, newUserRequest models.NewUserRequest) (fields models.UserFieldsTaken, err error)
	SetUserLastLogin(ctx context.Context, id int64) error

	CreateRefreshToken(ctx context.Context, refreshToken models.RefreshToken) (err error)
	GetRefreshToken(ctx context.Context, tokenHash string) (models.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
}

type Service struct {
	store                    Store
	accessTokenKey           []byte
	accessTokenTTL           time.Duration
	refreshTokenTTL          time.Duration
	refreshTokenCookieName   string
	refreshTokenCookieDomain string
}

type ServiceConfig struct {
	Store                    Store
	AccessTokenKey           string
	AccessTokenTTL           time.Duration
	RefreshTokenTTL          time.Duration
	RefreshTokenCookieName   string
	RefreshTokenCookieDomain string
}

func NewService(cfg ServiceConfig) *Service {
	return &Service{
		store:                    cfg.Store,
		accessTokenKey:           []byte(cfg.AccessTokenKey),
		accessTokenTTL:           cfg.AccessTokenTTL,
		refreshTokenTTL:          cfg.RefreshTokenTTL,
		refreshTokenCookieName:   cfg.RefreshTokenCookieName,
		refreshTokenCookieDomain: cfg.RefreshTokenCookieDomain,
	}
}

// jwtClaims defines the claims to be stored in the JWT
type jwtClaims struct {
	jwt.RegisteredClaims
	UserID  int64  `json:"user_id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
}

func (c jwtClaims) toModel() models.Claims {
	return models.Claims{
		UserID:  c.UserID,
		Email:   c.Email,
		IsAdmin: c.IsAdmin,
	}
}

// Register creates a user and starts an authenticated session
func (s *Service) Register(ctx context.Context, newUserRequest models.NewUserRequest, userAgent string, remember bool) (accessToken string, refreshToken string, err error) {
	// Validate new user
	if err = newUserRequest.Validate(); err != nil {
		return
	}

	taken, err := s.store.CheckTakenUserFields(ctx, newUserRequest)
	if err != nil {
		return
	}
	ve := models.ValidationErrors{}
	if taken.EmailTaken {
		ve.Add("email", "email already taken")
	}
	if taken.DisplayNameTaken {
		ve.Add("display_name", "display name already taken")
	}
	if ve.Any() {
		err = ve
		return
	}

	// Hash password
	hashedPassword, err := crypto.Hash(newUserRequest.Password)
	if err != nil {
		return
	}

	newUser := models.NewUser{
		Email:        newUserRequest.Email,
		DisplayName:  newUserRequest.DisplayName,
		Name:         newUserRequest.Name,
		PasswordHash: hashedPassword,
	}

	// Create user
	id, err := s.store.CreateUser(ctx, newUser)
	if err != nil {
		if errors.Is(err, store.ErrDuplicate) {
			err = ErrDuplicate
			return
		}
		return
	}

	// Fetch user and generate tokens
	user, err := s.store.GetUserByID(ctx, id)
	if err != nil {
		return
	}

	accessToken, err = s.GenerateAccessToken(user)
	if err != nil {
		return
	}

	if !remember {
		return
	}

	refreshToken, err = s.GenerateRefreshToken(ctx, user, userAgent)
	if err != nil {
		return
	}
	return
}

// Login authenticates a user and returns an access token and refresh token
func (s *Service) Login(ctx context.Context, email, password, userAgent string, remember bool) (accessToken string, refreshToken string, err error) {
	user, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		err = ErrInvalidCredentials
		return
	}

	if crypto.CompareHash(user.PasswordHash, password) != nil {
		err = ErrInvalidCredentials
		return
	}

	accessToken, err = s.GenerateAccessToken(user)
	if err != nil {
		return
	}

	if !remember {
		return
	}

	refreshToken, err = s.GenerateRefreshToken(ctx, user, userAgent)
	if err != nil {
		return
	}

	err = s.store.SetUserLastLogin(ctx, user.ID)
	if err != nil {
		return
	}

	return
}

// Logout invalidates a refresh token
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	refreshToken, err := crypto.Hash(refreshToken)
	if err != nil {
		return err
	}

	if err := s.store.DeleteRefreshToken(ctx, refreshToken); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil
		}
		return err
	}
	return nil
}

// Refresh refreshes an access token using a refresh token
func (s *Service) Refresh(ctx context.Context, refreshToken string) (newRefreshToken, accessToken string, err error) {
	// Hash provided refresh token
	refreshToken, err = crypto.Hash(refreshToken)
	if err != nil {
		return "", "", err
	}

	// Fetch refresh token
	token, err := s.store.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		err = ErrInvalidToken
		return "", "", err
	}

	// Revoke provided refresh token
	if err := s.store.DeleteRefreshToken(ctx, refreshToken); err != nil {
		return "", "", err
	}

	// Check for expired token
	if time.Now().UTC().After(token.ExpiresAt) {
		err = ErrInvalidToken
		return
	}

	// Fetch the user associated with the refresh token
	user, err := s.store.GetUserByID(ctx, token.UserID)
	if err != nil {
		err = ErrInvalidToken
		return
	}

	// Generate a new access token
	accessToken, err = s.GenerateAccessToken(user)
	if err != nil {
		return
	}

	// Rotate refresh tokens
	newRefreshToken, err = s.GenerateRefreshToken(ctx, user, token.UserAgent)
	if err != nil {
		return
	}
	return
}

// GenerateAccessToken generates an access token for a user
func (s *Service) GenerateAccessToken(user models.User) (string, error) {
	claims := &jwtClaims{
		UserID:  user.ID,
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(s.accessTokenTTL)),
			Issuer:    "critiquefi",
			Subject:   fmt.Sprint(user.ID),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(s.accessTokenKey)
}

// ValidateAccessToken parses and validates an access token, then returns the token claims
func (s *Service) ValidateAccessToken(accessToken string) (models.Claims, error) {
	t, err := jwt.ParseWithClaims(accessToken, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.accessTokenKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return models.Claims{}, ErrInvalidToken
	}

	claims, ok := t.Claims.(*jwtClaims)
	if !ok || !t.Valid {
		return models.Claims{}, ErrInvalidToken
	}

	claimsModel := claims.toModel()
	return claimsModel, nil
}

// GenerateRefreshToken generates a refresh token for a user
func (s *Service) GenerateRefreshToken(ctx context.Context, user models.User, userAgent string) (string, error) {
	refreshToken, err := crypto.RandomString(32)
	if err != nil {
		return "", err
	}

	hashedAccessToken, err := crypto.Hash(refreshToken)
	if err != nil {
		return "", err
	}

	token := models.RefreshToken{
		TokenHash: hashedAccessToken,
		UserID:    user.ID,
		UserAgent: userAgent,
		ExpiresAt: time.Now().Add(s.refreshTokenTTL).UTC(),
		CreatedAt: time.Now().UTC(),
	}

	if s.store.CreateRefreshToken(ctx, token) != nil {
		return "", err
	}

	return refreshToken, nil
}
