package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/andrew-hayworth22/critiquefi-service/internal/auth"
	"github.com/andrew-hayworth22/critiquefi-service/internal/models"
	"github.com/andrew-hayworth22/critiquefi-service/internal/store"
	"github.com/andrew-hayworth22/critiquefi-service/internal/testutil"
	"github.com/andrew-hayworth22/critiquefi-service/pkg/crypto"
)

func TestService_Register(t *testing.T) {
	t.Run("register", func(t *testing.T) {
		t.Parallel()
		tc := []struct {
			name                 string
			request              models.NewUserRequest
			userAgent            string
			remember             bool
			storeSetup           func(s *mockStore)
			expectedErr          error
			accessTokenExpected  bool
			refreshTokenExpected bool
		}{
			{
				name: "success: registered (don't remember)",
				request: models.NewUserRequest{
					Email:           "test.user@critiquefi.com",
					DisplayName:     "test.user",
					Name:            "Test User",
					Password:        "password",
					ConfirmPassword: "password",
				},
				userAgent: "test",
				remember:  false,
				storeSetup: func(s *mockStore) {
					s.on(CheckTakenUserFields, models.UserFieldsTaken{}, nil).
						on(CreateUser, int64(1), nil).
						on(GetUserByID, models.User{
							ID:           1,
							Email:        "user@critiquefi.com",
							DisplayName:  "user.name",
							Name:         "Test User",
							IsAdmin:      false,
							PasswordHash: "password",
							IsActive:     true,
						}, nil)
				},
				expectedErr:          nil,
				accessTokenExpected:  true,
				refreshTokenExpected: false,
			},
			{
				name: "success: registered (remember)",
				request: models.NewUserRequest{
					Email:           "test.user@critiquefi.com",
					DisplayName:     "test.user",
					Name:            "Test User",
					Password:        "password",
					ConfirmPassword: "password",
				},
				userAgent: "test",
				remember:  true,
				storeSetup: func(s *mockStore) {
					s.on(CheckTakenUserFields, models.UserFieldsTaken{}, nil).
						on(CreateUser, int64(1), nil).
						on(GetUserByID, models.User{
							ID:           1,
							Email:        "user@critiquefi.com",
							DisplayName:  "user.name",
							Name:         "Test User",
							IsAdmin:      false,
							PasswordHash: "password",
							IsActive:     true,
						}, nil).
						on(CreateRefreshToken, nil)
				},
				expectedErr:          nil,
				accessTokenExpected:  true,
				refreshTokenExpected: true,
			},
			{
				name: "error: validation errors - model validation",
				request: models.NewUserRequest{
					Email:           "critiquefi.com",
					DisplayName:     "er",
					Name:            "r",
					Password:        "ssword",
					ConfirmPassword: "password",
				},
				userAgent:  "test",
				remember:   true,
				storeSetup: func(s *mockStore) {},
				expectedErr: models.ValidationErrors{
					"email":            "invalid email address",
					"display_name":     "display name must be between 3 and 20 characters long",
					"name":             "name must be between 3 and 50 characters long",
					"password":         "password must be between 8 and 64 characters long",
					"confirm_password": "passwords do not match",
				},
				accessTokenExpected:  false,
				refreshTokenExpected: false,
			},
			{
				name: "error: validation errors - duplicate fields",
				request: models.NewUserRequest{
					Email:           "test.user@critiquefi.com",
					DisplayName:     "test.user",
					Name:            "Test User",
					Password:        "password",
					ConfirmPassword: "password",
				},
				userAgent: "test",
				remember:  true,
				storeSetup: func(s *mockStore) {
					s.on(CheckTakenUserFields, models.UserFieldsTaken{
						DisplayNameTaken: true,
						EmailTaken:       true,
					}, nil)
				},
				expectedErr: models.ValidationErrors{
					"email":        "email already taken",
					"display_name": "display name already taken",
				},
				accessTokenExpected:  false,
				refreshTokenExpected: false,
			},
			{
				name: "error: duplicate creation",
				request: models.NewUserRequest{
					Email:           "test.user@critiquefi.com",
					DisplayName:     "test.user",
					Name:            "Test User",
					Password:        "password",
					ConfirmPassword: "password",
				},
				userAgent: "test",
				remember:  true,
				storeSetup: func(s *mockStore) {
					s.on(CheckTakenUserFields, models.UserFieldsTaken{}, nil).
						on(CreateUser, int64(0), store.ErrDuplicate)
				},
				expectedErr:          auth.ErrDuplicate,
				accessTokenExpected:  false,
				refreshTokenExpected: false,
			},
		}

		for _, tc := range tc {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				s := newMockStore(t)
				tc.storeSetup(s)

				svc := auth.NewService(auth.ServiceConfig{
					Store:                    s,
					AccessTokenKey:           "key",
					AccessTokenTTL:           time.Hour,
					RefreshTokenTTL:          time.Hour,
					RefreshTokenCookieName:   "rt",
					RefreshTokenCookieDomain: "critiquefi.com",
				})

				accessToken, refreshToken, err := svc.Register(context.Background(), tc.request, tc.userAgent, tc.remember)
				testutil.CheckErr(err, tc.expectedErr, t)

				if tc.accessTokenExpected && accessToken == "" {
					t.Error("expected access token, got nothing")
				}
				if !tc.accessTokenExpected && accessToken != "" {
					t.Errorf("expected no access token, got %v", accessToken)
				}
				if tc.refreshTokenExpected && refreshToken == "" {
					t.Error("expected refresh token, got nothing")
				}
				if !tc.refreshTokenExpected && refreshToken != "" {
					t.Errorf("expected no refresh token, got %v", refreshToken)
				}
			})
		}
	})
}

func TestService_Login(t *testing.T) {
	t.Run("login", func(t *testing.T) {
		t.Parallel()

		hashedPassword, err := crypto.Hash("password")
		if err != nil {
			t.Fatal(err)
		}

		tc := []struct {
			name                 string
			email                string
			password             string
			userAgent            string
			remember             bool
			storeSetup           func(s *mockStore)
			expectedErr          error
			accessTokenExpected  bool
			refreshTokenExpected bool
		}{
			{
				name:      "success: logged in (don't remember)",
				email:     "user@critiquefi.com",
				password:  "password",
				userAgent: "test",
				remember:  false,
				storeSetup: func(s *mockStore) {
					s.on(GetUserByEmail, models.User{
						ID:           1,
						Email:        "user@critiquefi.com",
						DisplayName:  "test.user",
						Name:         "Test User",
						IsAdmin:      false,
						PasswordHash: hashedPassword,
						IsActive:     false,
					}, nil).on(SetUserLastLogin, nil)
				},
				expectedErr:          nil,
				accessTokenExpected:  true,
				refreshTokenExpected: false,
			},
			{
				name:      "success: logged in (remember)",
				email:     "user@critiquefi.com",
				password:  "password",
				userAgent: "test",
				remember:  true,
				storeSetup: func(s *mockStore) {
					s.on(GetUserByEmail, models.User{
						ID:           1,
						Email:        "user@critiquefi.com",
						DisplayName:  "test.user",
						Name:         "Test User",
						IsAdmin:      false,
						PasswordHash: hashedPassword,
						IsActive:     false,
					}, nil).on(SetUserLastLogin, nil).on(CreateRefreshToken, nil)
				},
				expectedErr:          nil,
				accessTokenExpected:  true,
				refreshTokenExpected: true,
			},
			{
				name:      "error: wrong email",
				email:     "nonexistent@critiquefi.com",
				password:  "password",
				userAgent: "test",
				remember:  false,
				storeSetup: func(s *mockStore) {
					s.on(GetUserByEmail, models.User{}, store.ErrNotFound)
				},
				expectedErr:          auth.ErrInvalidCredentials,
				accessTokenExpected:  false,
				refreshTokenExpected: false,
			},
			{
				name:      "error: wrong password",
				email:     "user@critiquefi.com",
				password:  "password",
				userAgent: "test",
				remember:  true,
				storeSetup: func(s *mockStore) {
					s.on(GetUserByEmail, models.User{
						ID:           1,
						Email:        "user@critiquefi.com",
						DisplayName:  "test.user",
						Name:         "Test User",
						IsAdmin:      false,
						PasswordHash: "wrongpassword",
						IsActive:     false,
					}, nil)
				},
				expectedErr:          auth.ErrInvalidCredentials,
				accessTokenExpected:  false,
				refreshTokenExpected: false,
			},
		}

		for _, tc := range tc {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				s := newMockStore(t)
				tc.storeSetup(s)

				svc := auth.NewService(auth.ServiceConfig{
					Store:                    s,
					AccessTokenKey:           "key",
					AccessTokenTTL:           time.Hour,
					RefreshTokenTTL:          time.Hour,
					RefreshTokenCookieName:   "rt",
					RefreshTokenCookieDomain: "critiquefi.com",
				})

				accessToken, refreshToken, err := svc.Login(context.Background(), tc.email, tc.password, tc.userAgent, tc.remember)
				testutil.CheckErr(err, tc.expectedErr, t)

				if tc.accessTokenExpected && accessToken == "" {
					t.Error("expected access token, got nothing")
				}
				if !tc.accessTokenExpected && accessToken != "" {
					t.Errorf("expected no access token, got %v", accessToken)
				}
				if tc.refreshTokenExpected && refreshToken == "" {
					t.Error("expected refresh token, got nothing")
				}
				if !tc.refreshTokenExpected && refreshToken != "" {
					t.Errorf("expected no refresh token, got %v", refreshToken)
				}
			})
		}
	})
}

func TestService_Logout(t *testing.T) {
	t.Run("logout", func(t *testing.T) {
		t.Parallel()
		tc := []struct {
			name         string
			refreshToken string
			storeSetup   func(s *mockStore)
			expectedErr  error
		}{
			{
				name:         "success: logged out",
				refreshToken: "test-refresh-token",
				storeSetup: func(s *mockStore) {
					s.on(DeleteRefreshToken, nil)
				},
				expectedErr: nil,
			},
			{
				name:         "success: no error if token not found",
				refreshToken: "test-refresh-token",
				storeSetup: func(s *mockStore) {
					s.on(DeleteRefreshToken, store.ErrNotFound)
				},
				expectedErr: nil,
			},
		}

		for _, tc := range tc {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				s := newMockStore(t)
				tc.storeSetup(s)

				svc := auth.NewService(auth.ServiceConfig{
					Store:                  s,
					AccessTokenKey:         "key",
					AccessTokenTTL:         time.Hour,
					RefreshTokenTTL:        time.Hour,
					RefreshTokenCookieName: "rt",
				})

				err := svc.Logout(context.Background(), tc.refreshToken)
				testutil.CheckErr(err, tc.expectedErr, t)
			})
		}
	})
}
